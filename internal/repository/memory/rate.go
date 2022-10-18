package memory

import "github.com/pkg/errors"

type RateM struct {
	rates map[string]float64
}

func NewRateM() *RateM {
	return &RateM{
		rates: make(map[string]float64),
	}
}

func (r *RateM) GetRate(code string) (float64, error) {
	rate, has := r.rates[code]
	if !has {
		return 0, errors.New("exchange rate not found by currency code")
	}
	return rate, nil
}

func (r *RateM) SaveRate(code string, rate float64) {
	r.rates[code] = rate
}
