package memory

import "github.com/pkg/errors"

type Currency struct {
	rates map[string]float64
}

func NewCurrency() *Currency {
	return &Currency{
		rates: make(map[string]float64),
	}
}

func (c *Currency) GetRate(code string) (float64, error) {
	rate, has := c.rates[code]
	if !has {
		return 0, errors.New("exchange rate not found by currency code")
	}
	return rate, nil
}

func (c *Currency) SaveRate(code string, rate float64) {
	c.rates[code] = rate
}
