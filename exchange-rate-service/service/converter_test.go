package service

import "testing"

func Test_Convert(t *testing.T) {

	converter := NewConverter()

	//mock rates
	rates := map[string]float64{
		"USD": 1.0,
		"INR": 83.12,
		"EUR": 0.92,
		"JPY": 149.50,
		"GBP": 0.79,
	}

	tests := []struct {
		name     string
		from     string
		to       string
		amount   float64
		expected float64
	}{
		{
			name:     "USD to INR",
			from:     "USD",
			to:       "INR",
			amount:   100,
			expected: 8312,
		},
		{
			name:     "EUR to GBP",
			from:     "EUR",
			to:       "GBP",
			amount:   100,
			expected: 85.87,
		},
		{
			name:     "Same currency",
			from:     "USD",
			to:       "USD",
			amount:   100,
			expected: 100,
		},
		{
			name:     "INR to EUR",
			from:     "INR",
			to:       "EUR",
			amount:   1000,
			expected: 11.07,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := converter.Convert(tt.from, tt.to, tt.amount, rates)

			diff := result - tt.expected
			if diff < -0.01 || diff > 0.01 {
				t.Errorf("Convert() = %v, expected %v", result, tt.expected)
			}
		})

	}
}
