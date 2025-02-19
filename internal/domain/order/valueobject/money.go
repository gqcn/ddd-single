package valueobject

import (
	"fmt"
	"math"
)

// Money represents a monetary value
type Money struct {
	amount   float64
	currency string
}

// NewMoney creates a new Money instance
func NewMoney(amount float64, currency string) *Money {
	return &Money{
		amount:   math.Round(amount*100) / 100, // Round to 2 decimal places
		currency: currency,
	}
}

// Amount returns the monetary amount
func (m *Money) Amount() float64 {
	return m.amount
}

// Currency returns the currency code
func (m *Money) Currency() string {
	return m.currency
}

// Add adds another Money value and returns a new Money instance
func (m *Money) Add(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("cannot add different currencies: %s and %s", m.currency, other.currency)
	}
	return NewMoney(m.amount+other.amount, m.currency), nil
}

// Multiply multiplies the amount by a factor and returns a new Money instance
func (m *Money) Multiply(factor float64) *Money {
	return NewMoney(m.amount*factor, m.currency)
}

// String returns the string representation of Money
func (m *Money) String() string {
	return fmt.Sprintf("%.2f %s", m.amount, m.currency)
}
