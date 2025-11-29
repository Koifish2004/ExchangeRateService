package service

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestConvert(t *testing.T) {
	converter := NewConverter()

	// Mock rates as decimals
	rates := map[string]decimal.Decimal{
		"USD": decimal.NewFromInt(1),
		"INR": decimal.NewFromFloat(83.12),
		"EUR": decimal.NewFromFloat(0.92),
		"JPY": decimal.NewFromFloat(149.50),
		"GBP": decimal.NewFromFloat(0.79),
	}

	tests := []struct {
		name     string
		from     string
		to       string
		amount   string
		expected string
	}{
		{
			name:     "USD to INR",
			from:     "USD",
			to:       "INR",
			amount:   "100",
			expected: "8312.00",
		},
		{
			name:     "EUR to GBP",
			from:     "EUR",
			to:       "GBP",
			amount:   "100",
			expected: "85.87",
		},
		{
			name:     "Same currency",
			from:     "USD",
			to:       "USD",
			amount:   "100",
			expected: "100.00",
		},
		{
			name:     "INR to EUR",
			from:     "INR",
			to:       "EUR",
			amount:   "1000",
			expected: "11.07",
		},
		{
			name:     "Decimal input",
			from:     "USD",
			to:       "INR",
			amount:   "100.50",
			expected: "8353.56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount := decimal.RequireFromString(tt.amount)
			result, err := converter.Convert(tt.from, tt.to, amount, rates)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Convert() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestConvertMissingRate(t *testing.T) {
	converter := NewConverter()

	rates := map[string]decimal.Decimal{
		"USD": decimal.NewFromInt(1),
		"INR": decimal.NewFromFloat(83.12),
	}

	amount := decimal.NewFromInt(100)
	_, err := converter.Convert("USD", "EUR", amount, rates)

	if err == nil {
		t.Error("Expected error for missing rate")
	}
}

func TestConvertZeroRate(t *testing.T) {
	converter := NewConverter()

	rates := map[string]decimal.Decimal{
		"USD": decimal.NewFromInt(1),
		"EUR": decimal.Zero,
	}

	amount := decimal.NewFromInt(100)
	_, err := converter.Convert("USD", "EUR", amount, rates)

	if err == nil {
		t.Error("Expected error for zero rate")
	}
}
