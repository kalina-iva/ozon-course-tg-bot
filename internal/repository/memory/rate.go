package memory

import "github.com/pkg/errors"

type Rate struct {
	rates map[string]float64
}

func NewRate() *Rate {
	return &Rate{
		rates: make(map[string]float64),
	}
}

func (c *Rate) GetRate(code string) (float64, error) {
	rate, has := c.rates[code]
	if !has {
		return 0, errors.New("exchange rate not found by currency code")
	}
	return rate, nil
}

func (c *Rate) SaveRate(code string, rate float64) {
	c.rates[code] = rate
}
