package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	appErrors "github.com/yourusername/exchange-rate-service/errors"
)

type LatestAPIResponse struct {
	Success bool                       `json:"success"`
	Error   map[string]string          `json:"error"`
	Quotes  map[string]decimal.Decimal `json:"quotes"`
}

type HistoricalAPIResponse struct {
	Success bool                       `json:"success"`
	Error   map[string]string          `json:"error"`
	Quotes  map[string]decimal.Decimal `json:"quotes"`
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

func (c *APIClient) FetchLatestRates() (map[string]decimal.Decimal, error) {

	u, _ := url.Parse(c.baseURL + "/live")
	q := u.Query()
	q.Set("access_key", c.apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, appErrors.APIFetchError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, appErrors.APIBadStatusError(resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, appErrors.APIResponseError(err)
	}

	var result LatestAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, appErrors.APIResponseError(err)
	}

	if !result.Success {
		errorMsg := "unknown error"
		if info, ok := result.Error["info"]; ok {
			errorMsg = info
		}
		return nil, appErrors.NewAPIError(errorMsg, nil)
	}

	normalized := make(map[string]decimal.Decimal)
	normalized["USD"] = decimal.NewFromInt(1)

	for pair, rate := range result.Quotes {
		if strings.HasPrefix(pair, "USD") {
			currency := pair[3:]
			normalized[currency] = rate
		}
	}

	return normalized, nil
}

func (c *APIClient) FetchHistoricalRates(date time.Time) (map[string]decimal.Decimal, error) {
	u, _ := url.Parse(c.baseURL + "/historical")
	q := u.Query()
	q.Set("access_key", c.apiKey)
	q.Set("date", date.Format("2006-01-02"))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())

	if err != nil {
		return nil, appErrors.APIFetchError(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, appErrors.APIBadStatusError(resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, appErrors.APIResponseError(err)
	}

	var result HistoricalAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, appErrors.APIResponseError(err)
	}

	if !result.Success {
		errorMsg := "unknown error"
		if info, ok := result.Error["info"]; ok {
			errorMsg = info
		}
		return nil, appErrors.NewAPIError(errorMsg, nil)
	}

	normalized := make(map[string]decimal.Decimal)
	normalized["USD"] = decimal.NewFromInt(1)
	for pair, rate := range result.Quotes {
		if strings.HasPrefix(pair, "USD") {
			currency := pair[3:]

			normalized[currency] = rate
		}
	}

	return normalized, nil
}
