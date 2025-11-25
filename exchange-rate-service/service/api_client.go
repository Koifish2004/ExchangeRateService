package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type LatestAPIResponse struct {
	Success bool               `json:"success"`
	Quotes  map[string]float64 `json:"quotes"`
}

type HistoricalAPIResponse struct {
	Success bool               `json:"success"`
	Quotes  map[string]float64 `json:"quotes"`
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
		baseURL: "https://api.exchangerate.host",
	}
}

func (c *APIClient) FetchLatestRates() (map[string]float64, error) {

	u, _ := url.Parse(c.baseURL + "/live")
	q := u.Query()
	q.Set("access_key", c.apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
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

	if !result.Success {
		return nil, fmt.Errorf("API returned error")
	}
	//can we only add the 5 currencies we are supporting?

	normalized := make(map[string]float64)
	normalized["USD"] = 1.0
	for pair, rate := range result.Quotes {
		if strings.HasPrefix(pair, "USD") {
			normalized[pair[3:]] = rate
		}
	}

	return normalized, nil
}

func (c *APIClient) FetchHistoricalRates(date time.Time) (map[string]float64, error) {
	u, _ := url.Parse(c.baseURL + "/historical")
	q := u.Query()
	q.Set("access_key", c.apiKey)
	q.Set("date", date.Format("2006-01-02"))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())

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

	if !result.Success {
		return nil, fmt.Errorf("API returned error")
	}
	//can we only add the 5 currencies we are supporting?
	normalized := make(map[string]float64)
	normalized["USD"] = 1.0
	for pair, rate := range result.Quotes {
		if strings.HasPrefix(pair, "USD") {
			normalized[pair[3:]] = rate
		}
	}

	return normalized, nil
}
