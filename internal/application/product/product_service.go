package product

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"main/internal/domain/product/entity"
	"main/internal/domain/product/service"
	"main/internal/domain/product/valueobject"
	sharedvo "main/internal/domain/shared/valueobject"
)

// ProductApplicationService 商品应用服务
// 应用服务负责用例编排，但不包含业务规则
type ProductApplicationService struct {
	productService *service.ProductService // 商品领域服务
}

// NewProductApplicationService 创建商品应用服务实例
func NewProductApplicationService(
	productService *service.ProductService,
) *ProductApplicationService {
	return &ProductApplicationService{
		productService: productService,
	}
}

// CreateProductCommand 创建商品命令
type CreateProductCommand struct {
	Name        string
	Description string
	Price       float64
	Stock       int
}

// CreateProduct 创建商品
// 应用服务方法负责：
// 1. 参数验证和转换
// 2. 用例编排
// 3. 事务处理
// 4. 不包含业务规则
func (s *ProductApplicationService) CreateProduct(ctx context.Context, cmd CreateProductCommand) (*entity.Product, error) {
	// 1. 转换命令到领域对象参数
	price := sharedvo.NewMoney(cmd.Price, "CNY")

	// 2. 调用领域服务创建商品
	product, err := s.productService.CreateProduct(
		ctx,
		"", // ID will be assigned by infrastructure layer
		cmd.Name,
		cmd.Description,
		price,
		cmd.Stock,
	)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to create product")
	}

	return product, nil
}

// UpdateProductCommand 更新商品命令
type UpdateProductCommand struct {
	Id          string
	Name        string
	Description string
	Price       float64
	Stock       int
	Status      valueobject.ProductStatus
}

// UpdateProduct 更新商品
func (s *ProductApplicationService) UpdateProduct(ctx context.Context, cmd UpdateProductCommand) (*entity.Product, error) {
	// 1. 转换命令到领域对象参数
	price := sharedvo.NewMoney(cmd.Price, "CNY")

	// 2. 调用领域服务更新商品
	product, err := s.productService.UpdateProduct(
		ctx,
		cmd.Id,
		cmd.Name,
		cmd.Description,
		price,
		cmd.Stock,
		cmd.Status,
	)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to update product")
	}

	return product, nil
}

// DeleteProductCommand 删除商品命令
type DeleteProductCommand struct {
	Id string
}

// DeleteProduct 删除商品
func (s *ProductApplicationService) DeleteProduct(ctx context.Context, cmd DeleteProductCommand) error {
	if err := s.productService.DeleteProduct(ctx, cmd.Id); err != nil {
		return gerror.Wrap(err, "failed to delete product")
	}
	return nil
}

// GetProductQuery 获取商品查询
type GetProductQuery struct {
	Id string
}

// GetProduct 获取商品
func (s *ProductApplicationService) GetProduct(ctx context.Context, query GetProductQuery) (*entity.Product, error) {
	product, err := s.productService.GetProduct(ctx, query.Id)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to get product")
	}
	return product, nil
}

// ListProductsQuery 列出商品查询
type ListProductsQuery struct {
	Status valueobject.ProductStatus
}

// ListProducts 列出商品
func (s *ProductApplicationService) ListProducts(ctx context.Context, query ListProductsQuery) ([]*entity.Product, error) {
	products, err := s.productService.ListProducts(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "failed to list products")
	}

	// 如果指定了状态，过滤商品列表
	if query.Status != "" {
		filtered := make([]*entity.Product, 0)
		for _, p := range products {
			if p.Status == query.Status {
				filtered = append(filtered, p)
			}
		}
		return filtered, nil
	}

	return products, nil
}

// ReserveStockCommand 预留库存命令
type ReserveStockCommand struct {
	ProductId string
	Quantity  int
}

// ReserveStock 预留库存
func (s *ProductApplicationService) ReserveStock(ctx context.Context, cmd ReserveStockCommand) error {
	if err := s.productService.ReserveStock(ctx, cmd.ProductId, cmd.Quantity); err != nil {
		return gerror.Wrap(err, "failed to reserve stock")
	}
	return nil
}

// ReleaseStockCommand 释放库存命令
type ReleaseStockCommand struct {
	ProductId string
	Quantity  int
}

// ReleaseStock 释放库存
func (s *ProductApplicationService) ReleaseStock(ctx context.Context, cmd ReleaseStockCommand) error {
	if err := s.productService.ReleaseStock(ctx, cmd.ProductId, cmd.Quantity); err != nil {
		return gerror.Wrap(err, "failed to release stock")
	}
	return nil
}
