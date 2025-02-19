package valueobject

import "github.com/gogf/gf/v2/errors/gerror"

var (
	ErrInvalidCurrency   = gerror.New("invalid currency")
	ErrCurrencyMismatch  = gerror.New("currency mismatch")
	ErrInvalidAmount     = gerror.New("invalid amount")
	ErrInvalidMultiplier = gerror.New("invalid multiplier")
)

// Money 金额值对象
// 在领域驱动设计中，Money 是一个典型的值对象：
// 1. 不可变性：所有操作都返回新的实例
// 2. 无副作用：不改变原有对象的状态
// 3. 完整性：包含了金额和货币单位
// 4. 自封含：包含了所有必要的业务规则
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
// 确保两个金额的货币单位相同
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
// 确保两个金额的货币单位相同
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
// 确保乘数为非负数
func (m *Money) Multiply(multiplier float64) *Money {
	return NewMoney(m.amount*multiplier, m.currency)
}

// Equals 判断金额是否相等
// 需要同时比较金额和货币单位
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
// 确保金额和货币单位的有效性
func (m *Money) Validate() error {
	if m.currency == "" {
		return gerror.Wrap(ErrInvalidCurrency, "currency cannot be empty")
	}
	if m.amount < 0 {
		return gerror.Wrap(ErrInvalidAmount, "amount cannot be negative")
	}
	return nil
}
