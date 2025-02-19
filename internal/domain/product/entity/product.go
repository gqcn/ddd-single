package entity

import (
	"main/internal/domain/product/valueobject"
	sharedvo "main/internal/domain/shared/valueobject"
)

// Product 商品实体
type Product struct {
	Id          string
	Name        string
	Description string
	Price       *sharedvo.Money
	Stock       int
	Status      valueobject.ProductStatus
	CreatedAt   int64
	UpdatedAt   int64
}

// NewProduct 创建商品实体
func NewProduct(
	id string,
	name string,
	description string,
	price *sharedvo.Money,
	stock int,
	status valueobject.ProductStatus,
	createdAt int64,
	updatedAt int64,
) *Product {
	return &Product{
		Id:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		Status:      status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// UpdateStock 更新库存
func (p *Product) UpdateStock(stock int) error {
	if stock < 0 {
		return valueobject.ErrInvalidStock
	}
	p.Stock = stock
	return nil
}

// UpdatePrice 更新价格
func (p *Product) UpdatePrice(price *sharedvo.Money) error {
	if price == nil {
		return valueobject.ErrInvalidPrice
	}
	p.Price = price
	return nil
}

// UpdateStatus 更新商品状态
func (p *Product) UpdateStatus(status valueobject.ProductStatus) error {
	if !status.IsValid() {
		return valueobject.ErrInvalidStatus
	}
	p.Status = status
	return nil
}

// Validate 验证商品
func (p *Product) Validate() error {
	if p.Id == "" {
		return valueobject.ErrInvalidId
	}
	if p.Name == "" {
		return valueobject.ErrInvalidName
	}
	if p.Price == nil {
		return valueobject.ErrInvalidPrice
	}
	if p.Stock < 0 {
		return valueobject.ErrInvalidStock
	}
	if !p.Status.IsValid() {
		return valueobject.ErrInvalidStatus
	}
	return nil
}
