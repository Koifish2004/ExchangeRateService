package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type LatestAPIResponse struct {
	Result          string             `json:"result"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

type HistoricalAPIResponse struct {
	Result          string             `json:"result"`
	Year            int                `json:"year"`
	Month           int                `json:"month"`
	Day             int                `json:"day"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

type APIClient struct {
	apiKey  string
	baseURL string
}

func NewClient() *APIClient {

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY not set in environment")
	}
	return &APIClient{
		apiKey:  apiKey,
		baseURL: "https://v6.exchangerate-api.com/v6",
	}
}

func (c *APIClient) FetchLatestRates() (map[string]float64, error) {

	url := fmt.Sprintf("%s/%s/latest/USD", c.baseURL, c.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	var result LatestAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if result.Result != "success" {
		return nil, fmt.Errorf("API returned error")
	}

	return result.ConversionRates, nil
}

func (c *APIClient) FetchHistoricalRates(date time.Time) (map[string]float64, error) {
	year := date.Year()
	month := date.Month()
	day := date.Day()

	url := fmt.Sprintf("%s/%s/history/USD/%d/%d/%d", c.baseURL, c.apiKey, year, month, day)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical data %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result HistoricalAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if result.Result != "success" {
		return nil, fmt.Errorf("API returned error")
	}

	return result.ConversionRates, nil
}
