package repository

import (
	"context"

	"main/internal/domain/order/entity"
)

// OrderRepository defines the interface for order persistence
type OrderRepository interface {
	// Save persists the order to the storage.
	Save(ctx context.Context, order *entity.Order) error

	// FindById retrieves an order by its id.
	FindById(ctx context.Context, id string) (*entity.Order, error)

	// FindByUserId retrieves all orders for a user.
	FindByUserId(ctx context.Context, userId string) ([]*entity.Order, error)

	// Update updates an existing order.
	Update(ctx context.Context, order *entity.Order) error

	// Delete removes an order from storage.
	Delete(ctx context.Context, id string) error
}
