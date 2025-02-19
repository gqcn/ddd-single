package order

import (
	"context"
	"time"

	"main/internal/domain/order/entity"
	"main/internal/domain/order/event"
	"main/internal/domain/order/repository"
	"main/internal/domain/order/service"
	"main/internal/domain/order/valueobject"
	"main/internal/infrastructure/eventbus"

	"github.com/gogf/gf/v2/errors/gerror"
)

// ApplicationService handles application-level order operations
type ApplicationService struct {
	orderRepo    repository.OrderRepository
	orderService *service.OrderService
	eventBus     eventbus.EventBus // 假设我们有一个事件总线接口
}

// NewApplicationService creates a new ApplicationService instance
func NewApplicationService(
	orderRepo repository.OrderRepository,
	orderService *service.OrderService,
	eventBus eventbus.EventBus,
) *ApplicationService {
	return &ApplicationService{
		orderRepo:    orderRepo,
		orderService: orderService,
		eventBus:     eventBus,
	}
}

// CreateOrder creates a new order
func (s *ApplicationService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*OrderDTO, error) {
	// Convert request to entity
	order := req.ToEntity()

	// Validate order
	if err := order.Validate(); err != nil {
		return nil, gerror.Wrap(err, "invalid order")
	}

	// Create order using domain service
	order, err := s.orderService.CreateOrder(ctx, order.UserId, order.Items)
	if err != nil {
		return nil, err
	}

	// Publish event
	s.eventBus.Publish(ctx, event.NewOrderCreatedEvent(order.Id, order.UserId))

	// Convert to DTO and return
	return ToDTO(order), nil
}

// GetOrder retrieves an order by Id
func (s *ApplicationService) GetOrder(ctx context.Context, orderId string) (*OrderDTO, error) {
	order, err := s.orderRepo.FindById(ctx, orderId)
	if err != nil {
		return nil, err
	}
	return ToDTO(order), nil
}

// GetUserOrders retrieves all orders for a user
func (s *ApplicationService) GetUserOrders(ctx context.Context, userId string) ([]*OrderDTO, error) {
	orders, err := s.orderRepo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return ToDTOs(orders), nil
}

// UpdateOrderStatus updates the status of an order
func (s *ApplicationService) UpdateOrderStatus(ctx context.Context, orderId string, req *UpdateOrderStatusRequest) error {
	// Get order
	order, err := s.orderRepo.FindById(ctx, orderId)
	if err != nil {
		return err
	}

	// Convert and validate new status
	newStatus, err := ToOrderStatus(req.Status)
	if err != nil {
		return err
	}

	// Store old status for event
	oldStatus := order.Status

	// Update status
	if err := order.UpdateStatus(newStatus); err != nil {
		return err
	}

	// Update order
	if err := s.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	// Publish event
	s.eventBus.Publish(event.NewOrderStatusChangedEvent(order.Id, oldStatus, newStatus, req.Remark))

	return nil
}

// CancelOrder cancels an order
func (s *ApplicationService) CancelOrder(ctx context.Context, orderId string, reason string) error {
	order, err := s.orderRepo.FindById(ctx, orderId)
	if err != nil {
		return err
	}

	oldStatus := order.Status

	if err := order.UpdateStatus(valueobject.OrderStatusCancelled); err != nil {
		return err
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	s.eventBus.Publish(event.NewOrderStatusChangedEvent(
		order.Id,
		oldStatus,
		valueobject.OrderStatusCancelled,
		reason,
	))

	return nil
}

// AddOrderItem adds an item to an existing order
func (s *ApplicationService) AddOrderItem(ctx context.Context, orderId string, itemReq *OrderItemRequest) error {
	order, err := s.orderRepo.FindById(ctx, orderId)
	if err != nil {
		return err
	}

	item := entity.NewOrderItem(
		itemReq.ProductId,
		itemReq.ProductName,
		itemReq.Quantity,
		itemReq.Price,
	)

	if err := order.AddItem(item); err != nil {
		return err
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	s.eventBus.Publish(&event.OrderItemAddedEvent{
		OrderEvent: event.OrderEvent{
			Id:        item.Id,
			OrderId:   orderId,
			Timestamp: time.Now(),
		},
		ProductId:   item.ProductId,
		ProductName: item.ProductName,
		Quantity:    item.Quantity,
		Price:       item.Price.Amount(),
	})

	return nil
}
