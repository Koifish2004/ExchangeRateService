package service

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type Cache struct {
	latestRates map[string]decimal.Decimal
	lastUpdated time.Time

	historicalRates map[string]map[string]decimal.Decimal

	mu sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		latestRates:     make(map[string]decimal.Decimal),
		historicalRates: make(map[string]map[string]decimal.Decimal),
	}
}

func (c *Cache) GetLatestRates() (map[string]decimal.Decimal, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.latestRates) == 0 {
		return nil, false
	}

	rates := make(map[string]decimal.Decimal)
	for i, j := range c.latestRates {
		rates[i] = j
	}

	return rates, true

}

func (c *Cache) SetLatestRates(rates map[string]decimal.Decimal) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.latestRates = rates
	c.lastUpdated = time.Now()
}

func (c *Cache) GetLastUpdated() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.lastUpdated
}

func (c *Cache) GetHistoricalRates(date time.Time) (map[string]decimal.Decimal, bool) {
	dateKey := date.Format("2006-01-02")
	c.mu.RLock()
	defer c.mu.RUnlock()

	rates, exists := c.historicalRates[dateKey]
	if !exists {
		return nil, false
	}

	ratesCopy := make(map[string]decimal.Decimal)
	for k, v := range rates {
		ratesCopy[k] = v
	}

	return ratesCopy, true

}

func (c *Cache) SetHistoricalRates(date time.Time, rates map[string]decimal.Decimal) {
	dateKey := date.Format("2006-01-02")
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.historicalRates[dateKey] = rates

}

func (c *Cache) ClearOldHistoricalData() {
	c.mu.Lock()
	defer c.mu.Unlock()

	cutoffDate := time.Now().AddDate(0, 0, -90)

	for dateStr := range c.historicalRates {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if date.Before(cutoffDate) {
			delete(c.historicalRates, dateStr)
		}
	}
}
