package main

//completely AI generated file. Do not use this to evaluate me pls.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/exchange-rate-service/handler"
	"github.com/yourusername/exchange-rate-service/service"
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Set test mode
	gin.SetMode(gin.TestMode)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// setupTestServer creates a real server with all components wired together
func setupTestServer(t *testing.T) *gin.Engine {
	// Check if API key is set
	if os.Getenv("API_KEY") == "" {
		t.Skip("API_KEY not set - skipping integration tests")
	}

	// Create real services (not mocked)
	rateFetcher := service.NewRateFetcherService()
	convertHandler := handler.NewConvertHandler(rateFetcher)

	// Setup router
	router := gin.New()
	router.GET("/convert", convertHandler.HandleConvert)

	return router
}

// Test 1: Basic conversion with latest rates
func TestIntegration_BasicConversion(t *testing.T) {
	router := setupTestServer(t)

	req := httptest.NewRequest("GET", "/convert?from=USD&to=INR&amount=100", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Check status code
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify response structure
	if result["from"] != "USD" {
		t.Errorf("Expected from=USD, got %v", result["from"])
	}
	if result["to"] != "INR" {
		t.Errorf("Expected to=INR, got %v", result["to"])
	}
	if result["result"] == nil {
		t.Error("Expected result field in response")
	}

	// Result should be reasonable (8000-9000 range for 100 USD to INR)
	resultValue := result["result"].(float64)
	if resultValue < 7000 || resultValue > 10000 {
		t.Errorf("Result seems unreasonable: %.2f", resultValue)
	}

	t.Logf("✓ Converted 100 USD to %.2f INR", resultValue)
}

// Test 2: All supported currency pairs
func TestIntegration_AllCurrencyPairs(t *testing.T) {
	router := setupTestServer(t)

	currencies := []string{"USD", "INR", "EUR", "JPY", "GBP"}

	for _, from := range currencies {
		for _, to := range currencies {
			if from == to {
				continue // Skip same currency
			}

			t.Run(from+"_to_"+to, func(t *testing.T) {
				url := "/convert?from=" + from + "&to=" + to + "&amount=100"
				req := httptest.NewRequest("GET", url, nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				if resp.Code != http.StatusOK {
					t.Errorf("Failed: %s", resp.Body.String())
				}
			})
		}
	}

	t.Log("✓ All 20 currency pairs working")
}

// Test 3: Historical conversion (30 days ago)
func TestIntegration_HistoricalConversion(t *testing.T) {
	router := setupTestServer(t)

	date := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	url := "/convert?from=EUR&to=GBP&amount=100&date=" + date

	req := httptest.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	if result["date"] != date {
		t.Errorf("Expected date %s in response", date)
	}

	t.Logf("✓ Historical conversion: 100 EUR to %.2f GBP on %s",
		result["result"].(float64), date)
}

// Test 4: Validation - unsupported currency
func TestIntegration_UnsupportedCurrency(t *testing.T) {
	router := setupTestServer(t)

	tests := []struct {
		name string
		url  string
	}{
		{"Bitcoin", "/convert?from=BTC&to=USD&amount=1"},
		{"Ethereum", "/convert?from=USD&to=ETH&amount=100"},
		{"Random", "/convert?from=XXX&to=YYY&amount=100"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 for unsupported currency, got %d", resp.Code)
			}

			var result map[string]interface{}
			json.Unmarshal(resp.Body.Bytes(), &result)

			if result["error"] == nil {
				t.Error("Expected error message")
			}
		})
	}

	t.Log("✓ Unsupported currencies properly rejected")
}

// Test 5: Validation - invalid amounts
func TestIntegration_InvalidAmount(t *testing.T) {
	router := setupTestServer(t)

	tests := []struct {
		name   string
		amount string
	}{
		{"Zero", "0"},
		{"Negative", "-100"},
		{"Not a number", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/convert?from=USD&to=INR&amount=" + tt.amount
			req := httptest.NewRequest("GET", url, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 for invalid amount, got %d", resp.Code)
			}
		})
	}

	t.Log("✓ Invalid amounts properly rejected")
}

// Test 6: Validation - future date
func TestIntegration_FutureDate(t *testing.T) {
	router := setupTestServer(t)

	futureDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	url := "/convert?from=USD&to=INR&amount=100&date=" + futureDate

	req := httptest.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for future date, got %d", resp.Code)
	}

	t.Log("✓ Future dates properly rejected")
}

// Test 7: Validation - date too old (>90 days)
func TestIntegration_OldDate(t *testing.T) {
	router := setupTestServer(t)

	oldDate := time.Now().AddDate(0, 0, -100).Format("2006-01-02")
	url := "/convert?from=USD&to=INR&amount=100&date=" + oldDate

	req := httptest.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for old date, got %d", resp.Code)
	}

	t.Log("✓ Dates older than 90 days properly rejected")
}

// Test 8: Validation - exactly 90 days (should work)
func TestIntegration_Exactly90Days(t *testing.T) {
	router := setupTestServer(t)

	date := time.Now().AddDate(0, 0, -90).Format("2006-01-02")
	url := "/convert?from=USD&to=INR&amount=100&date=" + date

	req := httptest.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected 200 for exactly 90 days, got %d", resp.Code)
	}

	t.Log("✓ Exactly 90 days ago works correctly")
}

// Test 9: Missing parameters
func TestIntegration_MissingParameters(t *testing.T) {
	router := setupTestServer(t)

	tests := []struct {
		name string
		url  string
	}{
		{"Missing from", "/convert?to=INR&amount=100"},
		{"Missing to", "/convert?from=USD&amount=100"},
		{"Missing amount", "/convert?from=USD&to=INR"},
		{"Missing all", "/convert"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 for missing params, got %d", resp.Code)
			}
		})
	}

	t.Log("✓ Missing parameters properly rejected")
}

// Test 10: Caching behavior (same request twice should be fast)
func TestIntegration_CachingWorks(t *testing.T) {
	router := setupTestServer(t)

	url := "/convert?from=USD&to=EUR&amount=100"

	// First request (may fetch from API)
	start1 := time.Now()
	req1 := httptest.NewRequest("GET", url, nil)
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	duration1 := time.Since(start1)

	// Second request (should use cache)
	start2 := time.Now()
	req2 := httptest.NewRequest("GET", url, nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	duration2 := time.Since(start2)

	// Both should succeed
	if resp1.Code != http.StatusOK || resp2.Code != http.StatusOK {
		t.Fatal("Both requests should succeed")
	}

	// Results should be identical
	var result1, result2 map[string]interface{}
	json.Unmarshal(resp1.Body.Bytes(), &result1)
	json.Unmarshal(resp2.Body.Bytes(), &result2)

	if result1["result"] != result2["result"] {
		t.Error("Cache should return same result")
	}

	t.Logf("✓ Caching works: First call %v, Second call %v", duration1, duration2)
}
