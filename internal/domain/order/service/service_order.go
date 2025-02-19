package service

import (
	"context"

	"main/internal/domain/order/entity"
	"main/internal/domain/order/event"
	"main/internal/domain/order/repository"
	"main/internal/domain/order/valueobject"
	"main/internal/infrastructure/eventbus"

	"github.com/gogf/gf/v2/errors/gerror"
)

// OrderService 领域服务，处理订单相关的核心业务逻辑
type OrderService struct {
	orderRepo repository.OrderRepository
	eventBus  eventbus.EventBus // 事件总线
}

// NewOrderService 创建订单领域服务实例
func NewOrderService(
	orderRepo repository.OrderRepository,
	eventBus eventbus.EventBus,
) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		eventBus:  eventBus,
	}
}

// CreateOrder 创建订单
// 这是一个领域服务方法，专注于订单领域的业务规则
func (s *OrderService) CreateOrder(ctx context.Context, userId string, items []*entity.OrderItem) (*entity.Order, error) {
	// 1. 创建订单实体
	order := entity.NewOrder(userId)

	// 2. 添加订单项
	for _, item := range items {
		if err := order.AddItem(item); err != nil {
			return nil, gerror.Wrap(err, "failed to add order item")
		}
	}

	// 3. 保存订单
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, gerror.Wrap(err, "failed to save order")
	}

	// 4. 发布订单创建事件
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

	// 2. 处理支付（调用领域实体的方法）
	if err := order.ProcessPayment(paymentInfo); err != nil {
		return gerror.Wrap(err, "failed to process payment")
	}

	// 3. 保存订单
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return gerror.Wrap(err, "failed to save order")
	}

	// 4. 发布订单支付事件
	if err := s.eventBus.Publish(ctx, event.NewOrderPaidEvent(order)); err != nil {
		return gerror.Wrap(err, "failed to publish order paid event")
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
		return gerror.New("order cannot be cancelled")
	}

	// 3. 取消订单
	if err = order.Cancel(); err != nil {
		return gerror.Wrap(err, "failed to cancel order")
	}

	// 4. 保存订单
	if err = s.orderRepo.Save(ctx, order); err != nil {
		return gerror.Wrap(err, "failed to save order")
	}

	// 5. 发布订单取消事件
	if err = s.eventBus.Publish(ctx, event.NewOrderCancelledEvent(order)); err != nil {
		return gerror.Wrap(err, "failed to publish order cancelled event")
	}

	return nil
}

// UpdateOrder 更新订单信息
// 这是一个领域服务方法，负责订单更新的持久化和事件发布
func (s *OrderService) UpdateOrder(ctx context.Context, order *entity.Order) error {
	// 1. 验证订单
	if err := order.Validate(); err != nil {
		return gerror.Wrap(err, "invalid order")
	}

	// 2. 保存订单
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return gerror.Wrap(err, "failed to update order")
	}

	// 3. 发布订单更新事件
	if err := s.eventBus.Publish(ctx, event.NewOrderUpdatedEvent(order)); err != nil {
		return gerror.Wrap(err, "failed to publish order updated event")
	}

	return nil
}

// GetOrder 获取订单
func (s *OrderService) GetOrder(ctx context.Context, orderId string) (*entity.Order, error) {
	order, err := s.orderRepo.FindById(ctx, orderId)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to find order")
	}
	return order, nil
}

// ListOrdersByUser 获取用户订单列表
func (s *OrderService) ListOrdersByUser(ctx context.Context, userId string, status valueobject.OrderStatus) ([]*entity.Order, error) {
	orders, err := s.orderRepo.FindByUserIdAndStatus(ctx, userId, status)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to find orders")
	}
	return orders, nil
}

// 内部辅助方法

// canCancelOrder 判断订单是否可以取消
func (s *OrderService) canCancelOrder(order *entity.Order) bool {
	return order.Status == valueobject.OrderStatusCreated
}
