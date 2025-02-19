package valueobject

import "github.com/gogf/gf/v2/errors/gerror"

var (
	ErrInvalidCurrency        = gerror.New("invalid currency")
	ErrCurrencyMismatch       = gerror.New("currency mismatch")
	ErrInvalidAmount          = gerror.New("invalid amount")
	ErrInvalidMultiplier      = gerror.New("invalid multiplier")
)

// Money 金额值对象
type Money struct {
	amount   float64
	currency string
}

// NewMoney 创建金额值对象
func NewMoney(amount float64, currency string) *Money {
	return &Money{
		amount:   amount,
		currency: currency,
	}
}

// Amount 获取金额
func (m *Money) Amount() float64 {
	return m.amount
}

// Currency 获取货币类型
func (m *Money) Currency() string {
	return m.currency
}

// Add 金额相加
func (m *Money) Add(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, gerror.Wrapf(ErrCurrencyMismatch, 
			"cannot add money with different currencies: %s and %s",
			m.currency,
			other.currency,
		)
	}
	return NewMoney(m.amount+other.amount, m.currency), nil
}

// Subtract 金额相减
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, gerror.Wrapf(ErrCurrencyMismatch,
			"cannot subtract money with different currencies: %s and %s",
			m.currency,
			other.currency,
		)
	}
	return NewMoney(m.amount-other.amount, m.currency), nil
}

// Multiply 金额乘以系数
func (m *Money) Multiply(multiplier float64) (*Money, error) {
	if multiplier < 0 {
		return nil, gerror.Wrap(ErrInvalidMultiplier, "multiplier cannot be negative")
	}
	return NewMoney(m.amount*multiplier, m.currency), nil
}

// Equals 判断金额是否相等
func (m *Money) Equals(other *Money) bool {
	if other == nil {
		return false
	}
	return m.amount == other.amount && m.currency == other.currency
}

// IsZero 判断金额是否为零
func (m *Money) IsZero() bool {
	return m.amount == 0
}

// IsNegative 判断金额是否为负数
func (m *Money) IsNegative() bool {
	return m.amount < 0
}

// IsPositive 判断金额是否为正数
func (m *Money) IsPositive() bool {
	return m.amount > 0
}

// Validate 验证金额
func (m *Money) Validate() error {
	if m.currency == "" {
		return gerror.Wrap(ErrInvalidCurrency, "currency cannot be empty")
	}
	if m.amount < 0 {
		return gerror.Wrap(ErrInvalidAmount, "amount cannot be negative")
	}
	return nil
}
