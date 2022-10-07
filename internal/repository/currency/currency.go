package currency

type Currency struct {
	rates map[string]float64
}

func New() *Currency {
	rates := make(map[string]float64)
	rates["RUB"] = 1
	rates["USD"] = 2
	rates["EUR"] = 3
	rates["CNY"] = 4
	return &Currency{
		rates: rates,
	}
}

func (c *Currency) GetRate(currency string) (float64, error) {
	return c.rates[currency], nil
}
