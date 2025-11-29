package main

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

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func setupTestServer(t *testing.T) *gin.Engine {
	if os.Getenv("API_KEY") == "" {
		t.Skip("API_KEY not set - skipping integration tests")
	}

	rateFetcher := service.NewRateFetcherService()
	convertHandler := handler.NewConvertHandler(rateFetcher)

	router := gin.New()
	router.GET("/convert", convertHandler.HandleConvert)

	return router
}

func TestIntegration_BasicConversion(t *testing.T) {
	router := setupTestServer(t)

	req := httptest.NewRequest("GET", "/convert?from=USD&to=INR&amount=100", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Logf("API Error Body: %s", resp.Body.String())

		if resp.Code == 502 || resp.Code == 429 {
			t.Skip("Skipping test due to API Rate Limit")
		}

		t.Fatalf("Expected 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Handler only returns {"amount": "result_value"}
	amountStr, ok := result["amount"].(string)
	if !ok {
		t.Errorf("Expected amount to be string, got %T", result["amount"])
	}

	t.Logf("✓ Converted 100 USD to INR: %s", amountStr)
}

func TestIntegration_DecimalAmount(t *testing.T) {
	router := setupTestServer(t)

	req := httptest.NewRequest("GET", "/convert?from=USD&to=EUR&amount=123.45", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	amountStr, ok := result["amount"].(string)
	if !ok {
		t.Errorf("Expected amount to be string, got %T", result["amount"])
	}

	t.Logf("✓ Converted 123.45 USD to EUR: %s", amountStr)
}

func TestIntegration_AllCurrencyPairs(t *testing.T) {
	router := setupTestServer(t)

	currencies := []string{"USD", "INR", "EUR", "JPY", "GBP"}

	for _, from := range currencies {
		for _, to := range currencies {
			if from == to {
				continue
			}

			t.Run(from+"_to_"+to, func(t *testing.T) {
				url := "/convert?from=" + from + "&to=" + to + "&amount=100"
				req := httptest.NewRequest("GET", url, nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				if resp.Code != http.StatusOK {
					t.Errorf("Failed: %s", resp.Body.String())
					return
				}

				var result map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &result)

				if _, ok := result["amount"].(string); !ok {
					t.Errorf("Expected string amount, got %T", result["amount"])
				}
			})
		}
	}

	t.Log("✓ All 20 currency pairs working")
}

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

	amountStr, ok := result["amount"].(string)
	if !ok {
		t.Fatalf("Expected string amount, got %T", result["amount"])
	}

	t.Logf("✓ Historical conversion: 100 EUR to GBP on %s = %s", date, amountStr)
}

func TestIntegration_UnsupportedCurrency(t *testing.T) {
	router := setupTestServer(t)

	tests := []struct {
		name string
		url  string
	}{
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
		})
	}

	t.Log("✓ Unsupported currencies properly rejected")
}

func TestIntegration_BTCSupported(t *testing.T) {
	router := setupTestServer(t)

	req := httptest.NewRequest("GET", "/convert?from=BTC&to=USD&amount=1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected 200 for BTC conversion, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	amountStr, ok := result["amount"].(string)
	if !ok {
		t.Errorf("Expected string amount, got %T", result["amount"])
	}

	t.Logf("✓ BTC supported: 1 BTC to USD = %s", amountStr)
}

func TestIntegration_InvalidAmount(t *testing.T) {
	router := setupTestServer(t)

	tests := []struct {
		name   string
		amount string
	}{
		{"Zero", "0"},
		{"Negative", "-100"},
		{"Not a number", "abc"},
		{"Invalid decimal", "100.5.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/convert?from=USD&to=INR&amount=" + tt.amount
			req := httptest.NewRequest("GET", url, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			if resp.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 for invalid amount '%s', got %d", tt.amount, resp.Code)
			}
		})
	}

	t.Log("✓ Invalid amounts properly rejected")
}

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

func TestIntegration_ResultConsistency(t *testing.T) {
	router := setupTestServer(t)

	url := "/convert?from=USD&to=EUR&amount=100"

	// Make two identical requests
	req1 := httptest.NewRequest("GET", url, nil)
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)

	req2 := httptest.NewRequest("GET", url, nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	var result1, result2 map[string]interface{}
	json.Unmarshal(resp1.Body.Bytes(), &result1)
	json.Unmarshal(resp2.Body.Bytes(), &result2)

	if result1["amount"] != result2["amount"] {
		t.Error("Cache should return identical results")
	}

	t.Logf("✓ Results consistent: %s", result1["amount"])
}

func TestIntegration_SameCurrencyConversion(t *testing.T) {
	router := setupTestServer(t)

	req := httptest.NewRequest("GET", "/convert?from=USD&to=USD&amount=100", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.Code)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &result)

	amountStr := result["amount"].(string)
	if amountStr != "100" {
		t.Errorf("Same currency conversion should return exact amount, got %s", amountStr)
	}

	t.Log("✓ Same currency conversion works correctly")
}
