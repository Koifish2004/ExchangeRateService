package service

import (
	appErrors "github.com/yourusername/exchange-rate-service/errors"
)

type Converter struct {
}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) Convert(from, to string, amount float64, rates map[string]float64) (float64, error) {
	fromRate, fromExists := rates[from]
	toRate, toExists := rates[to]

	if !fromExists {
		return 0, appErrors.MissingRateError(from)
	}

	if !toExists {
		return 0, appErrors.MissingRateError(to)
	}

	if fromRate == 0 {
		return 0, appErrors.InvalidRateError(from)
	}

	if toRate == 0 {
		return 0, appErrors.InvalidRateError(to)
	}

	amountInUSD := amount / fromRate
	result := amountInUSD * toRate

	return result, nil
}
