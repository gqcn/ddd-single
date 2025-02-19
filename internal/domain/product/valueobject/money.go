package valueobject

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
