package valueobject

import "github.com/gogf/gf/v2/errors/gerror"

// 订单领域错误定义
var (
	// ========================================================================
	// 支付相关错误
	// ========================================================================

	ErrInvalidPaymentAmount  = gerror.New("invalid payment amount")
	ErrInvalidPaymentMethod  = gerror.New("invalid payment method")
	ErrInvalidPaymentChannel = gerror.New("invalid payment channel")
	ErrInvalidPaymentTradeNo = gerror.New("invalid payment trade no")
	ErrInvalidPaymentTime    = gerror.New("invalid payment time")
	ErrPaymentAmountMismatch = gerror.New("payment amount does not match order total")
	ErrInvalidPaymentStatus  = gerror.New("invalid payment status")

	// ========================================================================
	// 订单状态相关错误
	// ========================================================================

	ErrInvalidOrderStatus      = gerror.New("invalid order status")
	ErrInvalidStatusTransition = gerror.New("invalid status transition")
	ErrOrderNotFound           = gerror.New("order not found")
	ErrOrderAlreadyExists      = gerror.New("order already exists")

	// ========================================================================
	// 订单项相关错误
	// ========================================================================

	ErrInvalidOrderItem  = gerror.New("invalid order item")
	ErrProductNotInOrder = gerror.New("product not in order")
	ErrCannotModifyOrder = gerror.New("cannot modify order in current status")
)
