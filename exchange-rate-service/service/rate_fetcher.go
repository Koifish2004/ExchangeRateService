package service

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	appErrors "github.com/yourusername/exchange-rate-service/errors"
)

var SupportedCurrencies = map[string]bool{
	"USD": true,
	"INR": true,
	"EUR": true,
	"JPY": true,
	"GBP": true,
	"BTC": true,
}

type RateFetcherService struct {
	apiClient *APIClient
	converter *Converter
	cache     *Cache
}

func NewRateFetcherService() *RateFetcherService {
	service := &RateFetcherService{
		apiClient: NewClient(),
		converter: NewConverter(),
		cache:     NewCache(),
	}

	service.loadLatestRates()

	return service
}

func (s *RateFetcherService) loadLatestRates() error {
	rates, err := s.apiClient.FetchLatestRates()

	if err != nil {
		return err
	}

	s.cache.SetLatestRates(rates)
	return nil

}

func (s *RateFetcherService) ConvertCurrency(from, to string, amount string, date *time.Time) (string, error) {
	err := s.validate(from, to, amount, date)
	if err != nil {
		return "", err
	}
	amountDecimal, _ := decimal.NewFromString(amount)

	var rates map[string]decimal.Decimal
	var err1 error

	if date == nil {
		rates, err1 = s.getLatestRates()
	} else {
		rates, err1 = s.getHistoricalRates(*date)
	}

	if err1 != nil {
		fmt.Printf("Error in rate-fetcher convert currency %v", err1)
		return "", err1
	}

	result, err := s.converter.Convert(from, to, amountDecimal, rates)
	if err != nil {
		return "", fmt.Errorf("conversion error: %v", err)
	}

	return result, nil
}

func (s *RateFetcherService) validate(from, to string, amountStr string, date *time.Time) error {
	if !SupportedCurrencies[from] {
		return appErrors.UnsupportedCurrencyError(from)
	}
	if !SupportedCurrencies[to] {
		return appErrors.UnsupportedCurrencyError(to)
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return appErrors.InvalidAmountError()
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return appErrors.InvalidAmountError()
	}

	if date != nil {

		now := time.Now().UTC()
		dateUTC := date.UTC()

		if dateUTC.Truncate(24 * time.Hour).After(now.Truncate(24 * time.Hour)) {
			return appErrors.FutureDateError()
		}

		cutoffDate := now.AddDate(0, 0, -90).Truncate(24 * time.Hour)
		if dateUTC.Truncate(24 * time.Hour).Before(cutoffDate) {
			return appErrors.DateTooOldError()
		}
	}

	return nil
}

func (s *RateFetcherService) getLatestRates() (map[string]decimal.Decimal, error) {

	rates, found := s.cache.GetLatestRates()
	if found {
		return rates, nil
	}

	return s.apiClient.FetchLatestRates()
}

func (s *RateFetcherService) getHistoricalRates(date time.Time) (map[string]decimal.Decimal, error) {
	rates, found := s.cache.GetHistoricalRates(date)

	if found {
		return rates, nil
	}

	ratescache, err := s.apiClient.FetchHistoricalRates(date)

	if err != nil {
		return nil, err
	}

	s.cache.SetHistoricalRates(date, ratescache)
	return ratescache, nil
}

func (s *RateFetcherService) StartHourlyRefresh() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			if err := s.loadLatestRates(); err != nil {
				fmt.Printf("Error refreshing rates: %v/n", err)
			} else {
				fmt.Println("Latest rates refreshed")
			}
			s.cache.ClearOldHistoricalData()
			fmt.Printf("Stats: %v/n", s.GetCacheStats())
		}
	}()

	fmt.Println("Hourly refresh started")

}

func (s *RateFetcherService) GetCacheStats() map[string]interface{} {
	lastUpdated := s.cache.GetLastUpdated()

	return map[string]interface{}{
		"last_updated":      lastUpdated.Format(time.RFC3339),
		"cache_age_minutes": time.Since(lastUpdated).Minutes(),
	}
}
