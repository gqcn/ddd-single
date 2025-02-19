package repository

import (
	"context"

	"main/internal/domain/order/entity"
	"main/internal/domain/order/valueobject"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// Save 保存订单
	Save(ctx context.Context, order *entity.Order) error

	// FindById 根据ID查找订单
	FindById(ctx context.Context, id string) (*entity.Order, error)

	// FindByUserIdAndStatus 根据用户ID和状态查找订单列表
	FindByUserIdAndStatus(ctx context.Context, userId string, status valueobject.OrderStatus) ([]*entity.Order, error)

	// Delete 删除订单
	Delete(ctx context.Context, id string) error

	// Update 更新订单
	Update(ctx context.Context, order *entity.Order) error
}
