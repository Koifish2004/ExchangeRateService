package service

import (
	"github.com/shopspring/decimal"
	appErrors "github.com/yourusername/exchange-rate-service/errors"
)

type Converter struct {
}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) Convert(from, to string, amount decimal.Decimal, rates map[string]decimal.Decimal) (string, error) {
	fromRate, fromExists := rates[from]
	toRate, toExists := rates[to]

	if !fromExists {
		return "", appErrors.MissingRateError(from)
	}

	if !toExists {
		return "", appErrors.MissingRateError(to)
	}

	if fromRate == decimal.NewFromInt(0) {
		return "", appErrors.InvalidRateError(from)
	}

	if toRate == decimal.NewFromInt(0) {
		return "", appErrors.InvalidRateError(to)
	}

	result := amount.Mul(toRate).Div(fromRate)

	return result.String(), nil
}
