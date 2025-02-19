package valueobject

import (
	"time"

	sharedvo "main/internal/domain/shared/valueobject"

	"github.com/gogf/gf/v2/errors/gerror"
)

// PaymentMethod 支付方法
type PaymentMethod string

const (
	PaymentMethodAlipay  PaymentMethod = "alipay"  // 支付宝
	PaymentMethodWechat  PaymentMethod = "wechat"  // 微信支付
	PaymentMethodBank    PaymentMethod = "bank"    // 银行卡
	PaymentMethodBalance PaymentMethod = "balance" // 余额支付
)

// PaymentChannel 支付渠道
type PaymentChannel string

const (
	PaymentChannelApp     PaymentChannel = "app"     // APP支付
	PaymentChannelH5      PaymentChannel = "h5"      // H5支付
	PaymentChannelWeb     PaymentChannel = "web"     // Web支付
	PaymentChannelQrCode  PaymentChannel = "qrcode"  // 二维码支付
	PaymentChannelCounter PaymentChannel = "counter" // 柜台支付
)

// PaymentInfo 支付信息值对象
type PaymentInfo struct {
	Amount         *sharedvo.Money  // 支付金额
	Method         PaymentMethod    // 支付方法
	Channel        PaymentChannel   // 支付渠道
	TradeNo        string          // 交易号
	ExtraData      interface{}     // 额外数据
	PaymentTime    int64           // 支付时间
}

// NewPaymentInfo 创建支付信息
func NewPaymentInfo(
	amount *sharedvo.Money,
	method PaymentMethod,
	channel PaymentChannel,
	tradeNo string,
	extraData interface{},
) *PaymentInfo {
	return &PaymentInfo{
		Amount:      amount,
		Method:      method,
		Channel:     channel,
		TradeNo:     tradeNo,
		ExtraData:   extraData,
		PaymentTime: time.Now().UnixMilli(),
	}
}

// Validate 验证支付信息
func (p *PaymentInfo) Validate() error {
	if p.Amount == nil || p.Amount.Amount() <= 0 {
		return gerror.New("invalid payment amount")
	}

	if !p.IsValidMethod() {
		return gerror.Newf("invalid payment method: %s", p.Method)
	}

	if !p.IsValidChannel() {
		return gerror.Newf("invalid payment channel: %s", p.Channel)
	}

	if p.TradeNo == "" {
		return gerror.New("trade number is required")
	}

	return nil
}

// IsValidMethod 验证支付方法是否有效
func (p *PaymentInfo) IsValidMethod() bool {
	switch p.Method {
	case PaymentMethodAlipay, PaymentMethodWechat,
		PaymentMethodBank, PaymentMethodBalance:
		return true
	default:
		return false
	}
}

// IsValidChannel 验证支付渠道是否有效
func (p *PaymentInfo) IsValidChannel() bool {
	switch p.Channel {
	case PaymentChannelApp, PaymentChannelH5,
		PaymentChannelWeb, PaymentChannelQrCode,
		PaymentChannelCounter:
		return true
	default:
		return false
	}
}
