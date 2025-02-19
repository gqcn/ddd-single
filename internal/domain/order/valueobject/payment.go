package valueobject

import "time"

// PaymentInfo 支付信息值对象
type PaymentInfo struct {
	Amount    *Money    // 支付金额
	Method    string    // 支付方式
	Channel   string    // 支付渠道
	TradeNo   string    // 交易号
	PaidAt    time.Time // 支付时间
	ExtraData map[string]interface{} // 额外数据
}

// NewPaymentInfo 创建支付信息值对象
func NewPaymentInfo(
	amount *Money,
	method string,
	channel string,
	tradeNo string,
	paidAt time.Time,
	extraData map[string]interface{},
) *PaymentInfo {
	return &PaymentInfo{
		Amount:    amount,
		Method:    method,
		Channel:   channel,
		TradeNo:   tradeNo,
		PaidAt:    paidAt,
		ExtraData: extraData,
	}
}

// Validate 验证支付信息
func (p *PaymentInfo) Validate() error {
	if p.Amount == nil {
		return ErrInvalidPaymentAmount
	}
	if p.Method == "" {
		return ErrInvalidPaymentMethod
	}
	if p.Channel == "" {
		return ErrInvalidPaymentChannel
	}
	if p.TradeNo == "" {
		return ErrInvalidPaymentTradeNo
	}
	if p.PaidAt.IsZero() {
		return ErrInvalidPaymentTime
	}
	return nil
}
