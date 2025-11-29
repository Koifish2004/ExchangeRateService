package service

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestLatestRates(t *testing.T) {
	cache := NewCache()

	rates, found := cache.GetLatestRates()

	if found {
		t.Error("Expected cache to be empty initially")
	}
	if rates != nil {
		t.Error("Expected nil when cache is empty")
	}

	testRates := map[string]decimal.Decimal{
		"USD": decimal.NewFromInt(1),
		"INR": decimal.NewFromFloat(83.12),
		"EUR": decimal.NewFromFloat(0.92),
	}

	cache.SetLatestRates(testRates)

	rates, found = cache.GetLatestRates()
	if !found {
		t.Error("Expected to find rates in cache")
	}

	expectedINR := decimal.NewFromFloat(83.12)
	if !rates["INR"].Equal(expectedINR) {
		t.Errorf("Expected INR rate %s, got %s", expectedINR, rates["INR"])
	}
}

func TestHistoricalRates(t *testing.T) {
	cache := NewCache()
	date := time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)

	rates, found := cache.GetHistoricalRates(date)
	if found {
		t.Error("Expected cache to be empty initially")
	}

	testRates := map[string]decimal.Decimal{
		"USD": decimal.NewFromInt(1),
		"INR": decimal.NewFromFloat(82.50),
	}

	cache.SetHistoricalRates(date, testRates)

	rates, found = cache.GetHistoricalRates(date)

	if !found {
		t.Error("Expected to find historical value")
	}

	expectedINR := decimal.NewFromFloat(82.50)
	if !rates["INR"].Equal(expectedINR) {
		t.Errorf("Expected INR rate %s, got %s", expectedINR, rates["INR"])
	}

	differentDate := time.Date(2025, 11, 2, 0, 0, 0, 0, time.UTC)
	_, found = cache.GetHistoricalRates(differentDate)
	if found {
		t.Error("Should not find rates for different date")
	}
}

func TestGetLastUpdated(t *testing.T) {
	cache := NewCache()

	lastUpdated := cache.GetLastUpdated()

	if !lastUpdated.IsZero() {
		t.Error("Initially, last update should be zero")
	}

	cache.SetLatestRates(map[string]decimal.Decimal{
		"USD": decimal.NewFromInt(1),
	})

	lastUpdated = cache.GetLastUpdated()

	if lastUpdated.IsZero() {
		t.Error("Last updated should be set after updating")
	}

	if time.Since(lastUpdated) > time.Second {
		t.Error("lastUpdated should be very recent")
	}
}

func TestClearHistoricalValues(t *testing.T) {
	cache := NewCache()
	oldDate := time.Now().AddDate(0, 0, -100)

	// Should be cleared
	cache.SetHistoricalRates(oldDate, map[string]decimal.Decimal{
		"USD": decimal.NewFromInt(1),
	})

	// Should remain
	recentDate := time.Now().AddDate(0, 0, -30)
	cache.SetHistoricalRates(recentDate, map[string]decimal.Decimal{
		"USD": decimal.NewFromInt(1),
	})

	cache.ClearOldHistoricalData()

	_, found := cache.GetHistoricalRates(oldDate)
	if found {
		t.Error("Old data should have been cleared")
	}

	_, found = cache.GetHistoricalRates(recentDate)
	if !found {
		t.Error("Recent data should not be cleared")
	}
}
