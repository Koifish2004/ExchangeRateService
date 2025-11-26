package service

import (
	"testing"
	"time"
)

func Test_LatestRates(t *testing.T) {
	cache := NewCache()

	rates, found := cache.GetLatestRates()

	if found {
		t.Error("Expected cache to be empty initially")
	}
	if rates != nil {
		t.Error("Expected nil when cache is empty")
	}

	testRates := map[string]float64{
		"USD": 1.0,
		"INR": 83.12,
		"EUR": 0.92,
	}

	cache.SetLatestRates(testRates)

	rates, found = cache.GetLatestRates()
	if !found {
		t.Error("Expected to find rates in cache")
	}

	if rates["INR"] != 83.12 {
		t.Errorf("Expected INR rate 83.12, got %f", rates["INR"])
	}

}

func Test_HistoricalRates(t *testing.T) {

	cache := NewCache()
	date := time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)

	rates, found := cache.GetHistoricalRates(date)
	if found {
		t.Error("Expected cache to be empty initially")
	}
	testRates := map[string]float64{
		"USD": 1.0,
		"INR": 82.50,
	}

	cache.SetHistoricalRates(date, testRates)

	rates, found = cache.GetHistoricalRates(date)

	if !found {
		t.Error("Expected to find historical value")
	}

	if rates["INR"] != 82.50 {
		t.Errorf("Expected INR rate 82.50, got %f", rates["INR"])
	}

	differentDate := time.Date(2025, 11, 2, 0, 0, 0, 0, time.UTC)
	_, found = cache.GetHistoricalRates(differentDate)
	if found {
		t.Error("Should not find rates for different date")
	}
}

func Test_GetLastUpdated(t *testing.T) {
	cache := NewCache()

	lastUpdated := cache.GetLastUpdated()

	if !lastUpdated.IsZero() {
		t.Error("Initially, last update should be zero")
	}

	cache.SetLatestRates(map[string]float64{
		"USD": 1.0,
	})

	lastUpdated = cache.GetLastUpdated()

	if lastUpdated.IsZero() {
		t.Error("Last updated should be set after updating")
	}

	if time.Since(lastUpdated) > time.Second {
		t.Error("lastUpdated should be very recent")
	}

}

func Test_ClearHistoricalValues(t *testing.T) {

	cache := NewCache()
	oldDate := time.Now().AddDate(0, 0, -100)

	//should be cleared
	cache.SetHistoricalRates(oldDate, map[string]float64{"USD": 1.0})

	//should remain
	recentDate := time.Now().AddDate(0, 0, -30)
	cache.SetHistoricalRates(recentDate, map[string]float64{"USD": 1.0})

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
