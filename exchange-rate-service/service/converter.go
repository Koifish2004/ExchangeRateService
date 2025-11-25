package service

import "fmt"

type Converter struct {
}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) Convert(from, to string, amount float64, rates map[string]float64) (float64, error) {
	fromRate, fromExists := rates[from]
	toRate, toExists := rates[to]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("missing exchange rate for currency: %s or %s", from, to)
	}

	if fromRate == 0 || toRate == 0 {
		return 0, fmt.Errorf("invalid exchange rate for currency: %s or %s", from, to)
	}

	amountInUSD := amount / fromRate
	result := amountInUSD * toRate

	return result, nil
}
