package order

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"main/internal/domain/order/entity"
	orderservice "main/internal/domain/order/service"
	"main/internal/domain/order/valueobject"
	productservice "main/internal/domain/product/service"
	sharedvo "main/internal/domain/shared/valueobject"
)

// OrderApplication 订单应用服务
// 应用服务负责用例编排和协调不同的领域服务
type OrderApplication struct {
	orderService   *orderservice.OrderService     // 订单领域服务
	productService *productservice.ProductService // 商品领域服务
}

// NewOrderApplication 创建订单应用服务实例
func NewOrderApplication(
	orderService *orderservice.OrderService,
	productService *productservice.ProductService,
) *OrderApplication {
	return &OrderApplication{
		orderService:   orderService,
		productService: productService,
	}
}

// CreateOrderCommand 创建订单命令
type CreateOrderCommand struct {
	UserId string
	Items  []OrderItemCommand
	Remark string
}

// OrderItemCommand 订单项命令
type OrderItemCommand struct {
	ProductId string
	Quantity  int
}

// CreateOrder 创建订单
// 应用服务方法负责：
// 1. 参数验证和转换
// 2. 协调不同领域服务
// 3. 事务处理
func (s *OrderApplication) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*entity.Order, error) {
	// 1. 验证商品信息并检查库存
	orderItems := make([]*entity.OrderItem, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		// 获取商品信息
		product, err := s.productService.GetProduct(ctx, item.ProductId)
		if err != nil {
			return nil, gerror.Wrap(err, "failed to get product")
		}

		// 检查库存
		if !s.productService.HasSufficientStock(ctx, item.ProductId, item.Quantity) {
			return nil, gerror.Newf(
				"insufficient stock for product %s",
				item.ProductId,
			)
		}

		// 创建订单项
		orderItem := entity.NewOrderItem(
			product.Id,
			product.Name,
			item.Quantity,
			product.Price.Amount(),
		)
		orderItems = append(orderItems, orderItem)
	}

	// 2. 创建订单（使用订单领域服务）
	order, err := s.orderService.CreateOrder(ctx, cmd.UserId, orderItems)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to create order")
	}

	// 3. 预扣库存
	for _, item := range cmd.Items {
		if err = s.productService.ReserveStock(ctx, item.ProductId, item.Quantity); err != nil {
			// 如果预扣库存失败，应该回滚订单创建
			// 这里可以通过发布事件来处理，或者使用分布式事务
			return nil, gerror.Wrap(err, "failed to reserve stock")
		}
	}

	// 4. 更新订单备注
	if cmd.Remark != "" {
		order.UpdateRemark(cmd.Remark)
		if err = s.orderService.UpdateOrder(ctx, order); err != nil {
			return nil, gerror.Wrap(err, "failed to update order remark")
		}
	}

	return order, nil
}

// PayOrderCommand 支付订单命令
type PayOrderCommand struct {
	OrderId        string
	Amount         float64
	PaymentMethod  valueobject.PaymentMethod
	PaymentChannel valueobject.PaymentChannel
	TradeNo        string
}

// PayOrder 支付订单
func (s *OrderApplication) PayOrder(ctx context.Context, cmd PayOrderCommand) error {
	// 1. 创建支付信息值对象
	paymentInfo := valueobject.NewPaymentInfo(
		sharedvo.NewMoney(cmd.Amount, "CNY"),
		cmd.PaymentMethod,
		cmd.PaymentChannel,
		cmd.TradeNo,
		nil,
	)

	// 2. 调用领域服务处理支付
	if err := s.orderService.PayOrder(ctx, cmd.OrderId, paymentInfo); err != nil {
		return gerror.Wrap(err, "failed to pay order")
	}

	return nil
}

// CancelOrderCommand 取消订单命令
type CancelOrderCommand struct {
	OrderId string
}

// CancelOrder 取消订单
func (s *OrderApplication) CancelOrder(ctx context.Context, cmd CancelOrderCommand) error {
	// 1. 获取订单信息
	order, err := s.orderService.GetOrder(ctx, cmd.OrderId)
	if err != nil {
		return gerror.Wrap(err, "failed to get order")
	}

	// 2. 调用领域服务取消订单
	if err := s.orderService.CancelOrder(ctx, cmd.OrderId); err != nil {
		return gerror.Wrap(err, "failed to cancel order")
	}

	// 3. 释放库存
	for _, item := range order.Items {
		if err := s.productService.ReleaseStock(ctx, item.ProductId, item.Quantity); err != nil {
			// 如果释放库存失败，应该通过事件或其他方式来处理不一致
			return gerror.Wrap(err, "failed to release stock")
		}
	}

	return nil
}

// GetOrderQuery 获取订单查询
type GetOrderQuery struct {
	OrderId string
}

// GetOrder 获取订单
func (s *OrderApplication) GetOrder(ctx context.Context, query GetOrderQuery) (*entity.Order, error) {
	order, err := s.orderService.GetOrder(ctx, query.OrderId)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to get order")
	}
	return order, nil
}

// ListOrdersByUserQuery 获取用户订单列表查询
type ListOrdersByUserQuery struct {
	UserId string
	Status valueobject.OrderStatus
}

// ListOrdersByUser 获取用户订单列表
func (s *OrderApplication) ListOrdersByUser(ctx context.Context, query ListOrdersByUserQuery) ([]*entity.Order, error) {
	orders, err := s.orderService.ListOrdersByUser(ctx, query.UserId, query.Status)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to list orders")
	}
	return orders, nil
}
