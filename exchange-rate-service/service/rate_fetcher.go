package service

import (
	"time"
)

var SupportedCurrencies = map[string]bool{
	"USD": true,
	"INR": true,
	"EUR": true,
	"JPY": true,
	"GBP": true,
}

type RateFetcherService struct {
	apiClient *APIClient
}

func (s *RateFetcherService) ConvertCurrency(from, to string, amount float64, date *time.Time) (float64, error) {
	return 10.43, nil
}
