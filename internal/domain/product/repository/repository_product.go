package repository

import (
	"context"

	"main/internal/domain/product/entity"
)

// ProductRepository 商品仓储接口
type ProductRepository interface {
	// Save 保存商品
	Save(ctx context.Context, product *entity.Product) error
	// FindById 根据Id查找商品
	FindById(ctx context.Context, id string) (*entity.Product, error)
	// FindAll 查找所有商品
	FindAll(ctx context.Context) ([]*entity.Product, error)
	// Update 更新商品
	Update(ctx context.Context, product *entity.Product) error
	// Delete 删除商品
	Delete(ctx context.Context, id string) error
}
