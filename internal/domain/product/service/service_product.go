package service

import (
	"context"

	"main/internal/domain/product/valueobject"

	"github.com/gogf/gf/v2/errors/gerror"

	"main/internal/domain/product/entity"
	"main/internal/domain/product/repository"
	sharedvo "main/internal/domain/shared/valueobject"
)

// ProductService 商品服务
type ProductService struct {
	productRepo repository.ProductRepository
}

// NewProductService 创建商品服务实例
func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// CreateProduct 创建商品
func (s *ProductService) CreateProduct(
	ctx context.Context,
	id string,
	name string,
	description string,
	price *sharedvo.Money,
	stock int,
) (*entity.Product, error) {
	// 检查商品是否已存在
	existing, err := s.productRepo.FindById(ctx, id)
	if err != nil && !gerror.Is(err, valueobject.ErrProductNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, valueobject.ErrProductExists
	}

	// 创建新商品
	product := entity.NewProduct(
		id,
		name,
		description,
		price,
		stock,
		valueobject.ProductStatusDraft,
		0, // createdAt will be set by repository
		0, // updatedAt will be set by repository
	)

	// 验证商品
	if err = product.Validate(); err != nil {
		return nil, err
	}

	// 保存商品
	if err = s.productRepo.Save(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProduct 更新商品
func (s *ProductService) UpdateProduct(
	ctx context.Context,
	id string,
	name string,
	description string,
	price *sharedvo.Money,
	stock int,
	status valueobject.ProductStatus,
) (*entity.Product, error) {
	// 获取现有商品
	product, err := s.productRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新商品信息
	product.Name = name
	product.Description = description
	if err = product.UpdatePrice(price); err != nil {
		return nil, err
	}
	if err = product.UpdateStock(stock); err != nil {
		return nil, err
	}
	if err = product.UpdateStatus(status); err != nil {
		return nil, err
	}

	// 验证商品
	if err = product.Validate(); err != nil {
		return nil, err
	}

	// 保存更新
	if err = s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct 获取商品
func (s *ProductService) GetProduct(ctx context.Context, id string) (*entity.Product, error) {
	return s.productRepo.FindById(ctx, id)
}

// ListProducts 获取所有商品
func (s *ProductService) ListProducts(ctx context.Context) ([]*entity.Product, error) {
	return s.productRepo.FindAll(ctx)
}

// DeleteProduct 删除商品
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	// 检查商品是否存在
	product, err := s.productRepo.FindById(ctx, id)
	if err != nil {
		return err
	}

	// 更新状态为已删除
	if err := product.UpdateStatus(valueobject.ProductStatusDeleted); err != nil {
		return err
	}

	// 保存更新
	return s.productRepo.Update(ctx, product)
}

// HasSufficientStock 检查商品是否有足够的库存
func (s *ProductService) HasSufficientStock(ctx context.Context, productId string, quantity int) bool {
	product, err := s.productRepo.FindById(ctx, productId)
	if err != nil {
		return false
	}

	return product.Stock >= quantity && product.Status == valueobject.ProductStatusOnSale
}

// ReserveStock 预扣商品库存
func (s *ProductService) ReserveStock(ctx context.Context, productId string, quantity int) error {
	product, err := s.productRepo.FindById(ctx, productId)
	if err != nil {
		return gerror.Wrap(err, "failed to find product")
	}

	// 检查库存是否充足
	if !s.HasSufficientStock(ctx, productId, quantity) {
		return gerror.Newf("insufficient stock for product %s", productId)
	}

	// 更新库存
	if err := product.UpdateStock(product.Stock - quantity); err != nil {
		return gerror.Wrap(err, "failed to update stock")
	}

	// 如果库存为0，更新状态为售罄
	if product.Stock == 0 {
		if err := product.UpdateStatus(valueobject.ProductStatusSoldOut); err != nil {
			return gerror.Wrap(err, "failed to update product status")
		}
	}

	// 保存更新
	return s.productRepo.Update(ctx, product)
}

// ReleaseStock 释放商品库存
func (s *ProductService) ReleaseStock(ctx context.Context, productId string, quantity int) error {
	product, err := s.productRepo.FindById(ctx, productId)
	if err != nil {
		return gerror.Wrap(err, "failed to find product")
	}

	// 更新库存
	if err := product.UpdateStock(product.Stock + quantity); err != nil {
		return gerror.Wrap(err, "failed to update stock")
	}

	// 如果商品之前是售罄状态，且现在有库存了，更新状态为在售
	if product.Status == valueobject.ProductStatusSoldOut && product.Stock > 0 {
		if err := product.UpdateStatus(valueobject.ProductStatusOnSale); err != nil {
			return gerror.Wrap(err, "failed to update product status")
		}
	}

	// 保存更新
	return s.productRepo.Update(ctx, product)
}
