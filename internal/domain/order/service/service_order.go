package service

import (
	"context"

	"main/internal/domain/order/entity"
	"main/internal/domain/order/event"
	"main/internal/domain/order/repository"
	"main/internal/domain/order/valueobject"
	"main/internal/domain/product/service"
	"main/internal/infrastructure/eventbus"

	"github.com/gogf/gf/v2/errors/gerror"
)

// OrderService 领域服务，处理订单相关的核心业务逻辑
type OrderService struct {
	orderRepo      repository.OrderRepository
	productService *service.ProductService // 商品领域服务
	eventBus       eventbus.EventBus       // 事件总线
}

// NewOrderService 创建订单领域服务实例
func NewOrderService(
	orderRepo repository.OrderRepository,
	productService *service.ProductService,
	eventBus eventbus.EventBus,
) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		productService: productService,
		eventBus:       eventBus,
	}
}

// CreateOrder 创建订单
// 这是一个领域服务方法，因为它需要协调多个实体（订单和商品）和确保业务规则
func (s *OrderService) CreateOrder(ctx context.Context, userId string, items []*entity.OrderItem) (*entity.Order, error) {
	// 1. 创建订单实体
	order := entity.NewOrder(userId)

	// 2. 验证商品信息并检查库存
	for _, item := range items {
		// 调用商品领域服务验证商品
		product, err := s.productService.GetProduct(ctx, item.ProductId)
		if err != nil {
			return nil, gerror.Wrap(err, "failed to get product")
		}

		// 验证商品价格
		if product.Price.Amount() != item.Price.Amount() {
			return nil, gerror.Newf(
				"product price mismatch: expected %v, got %v",
				product.Price.Amount(),
				item.Price.Amount(),
			)
		}

		// 检查库存
		if !s.productService.HasSufficientStock(ctx, item.ProductId, item.Quantity) {
			return nil, gerror.Newf(
				"insufficient stock for product %s",
				item.ProductId,
			)
		}

		// 添加订单项
		if err = order.AddItem(item); err != nil {
			return nil, err
		}
	}

	// 3. 保存订单
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, gerror.Wrap(err, "failed to save order")
	}

	// 4. 预扣库存（这里可能需要通过领域事件来实现最终一致性）
	if err := s.reserveStock(ctx, order); err != nil {
		return nil, gerror.Wrap(err, "failed to reserve stock")
	}

	// 5. 发布订单创建事件
	if err := s.eventBus.Publish(ctx, event.NewOrderCreatedEvent(order)); err != nil {
		return nil, gerror.Wrap(err, "failed to publish order created event")
	}

	return order, nil
}

// PayOrder 支付订单
func (s *OrderService) PayOrder(ctx context.Context, orderId string, paymentInfo *valueobject.PaymentInfo) error {
	// 1. 获取订单
	order, err := s.orderRepo.FindById(ctx, orderId)
	if err != nil {
		return gerror.Wrap(err, "failed to find order")
	}

	// 2. 验证订单状态
	if order.Status != valueobject.OrderStatusCreated {
		return gerror.Newf("invalid order status for payment: %s", order.Status)
	}

	// 3. 验证支付金额
	if paymentInfo.Amount.Amount() != order.TotalAmount.Amount() {
		return gerror.New("payment amount does not match order total")
	}

	// 4. 更新订单状态
	oldStatus := order.Status
	if err = order.UpdateStatus(valueobject.OrderStatusPaid); err != nil {
		return err
	}

	// 5. 保存订单
	if err = s.orderRepo.Update(ctx, order); err != nil {
		return gerror.Wrap(err, "failed to update order")
	}

	// 6. 发布订单状态变更事件
	if err = s.eventBus.Publish(ctx, event.NewOrderStatusChangedEvent(
		order.Id,
		string(oldStatus),
		string(order.Status),
	)); err != nil {
		return gerror.Wrap(err, "failed to publish order status changed event")
	}

	return nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, orderId string) error {
	// 1. 获取订单
	order, err := s.orderRepo.FindById(ctx, orderId)
	if err != nil {
		return gerror.Wrap(err, "failed to find order")
	}

	// 2. 验证订单是否可以取消
	if !s.canCancelOrder(order) {
		return gerror.Newf("order cannot be cancelled in status: %s", order.Status)
	}

	// 3. 更新订单状态
	oldStatus := order.Status
	if err = order.UpdateStatus(valueobject.OrderStatusCancelled); err != nil {
		return err
	}

	// 4. 释放库存
	if err = s.releaseStock(ctx, order); err != nil {
		return gerror.Wrap(err, "failed to release stock")
	}

	// 5. 保存订单
	if err = s.orderRepo.Update(ctx, order); err != nil {
		return gerror.Wrap(err, "failed to update order")
	}

	// 6. 发布订单取消事件
	if err = s.eventBus.Publish(ctx, event.NewOrderCanceledEvent(order.Id, "user requested")); err != nil {
		return gerror.Wrap(err, "failed to publish order canceled event")
	}

	// 7. 发布订单状态变更事件
	if err = s.eventBus.Publish(ctx, event.NewOrderStatusChangedEvent(
		order.Id,
		string(oldStatus),
		string(order.Status),
	)); err != nil {
		return gerror.Wrap(err, "failed to publish order status changed event")
	}

	return nil
}

// ValidateOrderStatus 验证订单状态转换
func (s *OrderService) ValidateOrderStatus(currentStatus, newStatus valueobject.OrderStatus) error {
	if !currentStatus.CanTransitionTo(newStatus) {
		return gerror.Newf(
			"invalid status transition from %s to %s",
			currentStatus,
			newStatus,
		)
	}
	return nil
}

// 内部辅助方法

// canCancelOrder 判断订单是否可以取消
func (s *OrderService) canCancelOrder(order *entity.Order) bool {
	return order.Status == valueobject.OrderStatusCreated ||
		order.Status == valueobject.OrderStatusPaid
}

// reserveStock 预扣库存
func (s *OrderService) reserveStock(ctx context.Context, order *entity.Order) error {
	for _, item := range order.Items {
		if err := s.productService.ReserveStock(ctx, item.ProductId, item.Quantity); err != nil {
			return err
		}
	}
	return nil
}

// releaseStock 释放库存
func (s *OrderService) releaseStock(ctx context.Context, order *entity.Order) error {
	for _, item := range order.Items {
		if err := s.productService.ReleaseStock(ctx, item.ProductId, item.Quantity); err != nil {
			return err
		}
	}
	return nil
}
