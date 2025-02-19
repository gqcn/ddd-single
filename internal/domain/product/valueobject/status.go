package valueobject

import "errors"

// 商品相关错误
var (
	ErrInvalidId          = errors.New("invalid product id")
	ErrInvalidName        = errors.New("invalid product name")
	ErrInvalidPrice       = errors.New("invalid product price")
	ErrInvalidStock       = errors.New("invalid product stock")
	ErrInvalidStatus      = errors.New("invalid product status")
	ErrProductNotFound    = errors.New("product not found")
	ErrProductExists      = errors.New("product already exists")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrProductUnavailable = errors.New("product is unavailable")
)

// ProductStatus 商品状态
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"     // 草稿
	ProductStatusOnSale    ProductStatus = "on_sale"   // 在售
	ProductStatusOffSale   ProductStatus = "off_sale"  // 下架
	ProductStatusSoldOut   ProductStatus = "sold_out"  // 售罄
	ProductStatusDeleted   ProductStatus = "deleted"   // 删除
)

// IsValid 检查状态是否有效
func (s ProductStatus) IsValid() bool {
	switch s {
	case ProductStatusDraft, ProductStatusOnSale, ProductStatusOffSale,
		ProductStatusSoldOut, ProductStatusDeleted:
		return true
	default:
		return false
	}
}

// CanTransitionTo 检查是否可以转换到目标状态
func (s ProductStatus) CanTransitionTo(target ProductStatus) bool {
	switch s {
	case ProductStatusDraft:
		return target == ProductStatusOnSale || target == ProductStatusDeleted
	case ProductStatusOnSale:
		return target == ProductStatusOffSale || target == ProductStatusSoldOut
	case ProductStatusOffSale:
		return target == ProductStatusOnSale || target == ProductStatusDeleted
	case ProductStatusSoldOut:
		return target == ProductStatusOnSale || target == ProductStatusDeleted
	case ProductStatusDeleted:
		return false
	default:
		return false
	}
}

// String 返回状态的字符串表示
func (s ProductStatus) String() string {
	return string(s)
}
