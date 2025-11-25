package service

import (
	"fmt"
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

func (s *RateFetcherService) ConvertCurrency(from, to string, amount float64, date *time.Time) (float64, error) {
	err := s.validate(from, to, amount, date)
	if err != nil {
		return 0, err
	}

	var rates map[string]float64
	var err1 error

	if date == nil {
		rates, err1 = s.getLatestRates()
	} else {
		rates, err1 = s.getHistoricalRates(*date)
	}

	if err1 != nil {
		fmt.Printf("Error in rate-fetcher convert currency %v", err1)
		return 0, err1
	}

	result, err := s.converter.Convert(from, to, amount, rates)
	if err != nil {
		return 0, fmt.Errorf("conversion error: %v", err)
	}

	return result, nil
}

func (s *RateFetcherService) validate(from, to string, amount float64, date *time.Time) error {
	if !SupportedCurrencies[from] {
		return fmt.Errorf("unsupported currency %s", from)
	}
	if !SupportedCurrencies[to] {
		fmt.Printf("Validation failed: unsupported currency %s\n", to)
		return fmt.Errorf("unsupported currency %s", to)
	}

	if amount <= 0 {
		return fmt.Errorf("amount cannot be 0 or negative")
	}

	if date != nil {
		if date.After(time.Now()) {
			return fmt.Errorf("date cannot be future")
		}

		if date.Before(time.Now().AddDate(0, 0, -90)) {
			return fmt.Errorf("max lookback is 90 days, date too old")
		}
	}

	return nil
}

func (s *RateFetcherService) getLatestRates() (map[string]float64, error) {

	rates, found := s.cache.GetLatestRates()
	if found {
		return rates, nil
	}

	return s.apiClient.FetchLatestRates()
}

func (s *RateFetcherService) getHistoricalRates(date time.Time) (map[string]float64, error) {
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
