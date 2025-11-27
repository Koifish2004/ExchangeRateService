# Exchange Rate Service

A high-performance currency exchange rate service built in Go that provides real-time and historical currency conversion with in-memory caching.

## Features

- Real-time currency conversion for USD, INR, EUR, JPY, GBP
- Historical exchange rates (up to 90 days)
- In-memory caching for optimal performance
- Hourly automatic rate refresh
- Thread-safe concurrent request handling
- Comprehensive input validation
- RESTful API design
- Docker support

## Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│  HTTP Handler   │  (Gin Framework)
└────────┬────────┘
         │
         ▼
┌──────────────────┐
│ Rate Fetcher     │  (Business Logic)
│   Service        │
└────┬───────┬─────┘
     │       │
     ▼       ▼
┌────────┐ ┌──────────┐
│ Cache  │ │   API    │
│ (RAM)  │ │  Client  │
└────────┘ └─────┬────┘
                 │
                 ▼
         ┌───────────────┐
         │ External API  │
         │(exchangerate- │
         │    api.com)   │
         └───────────────┘
```

### Components

1. **API Client** (`service/api_client.go`) - Handles external API calls
2. **Cache** (`service/cache.go`) - Thread-safe in-memory storage
3. **Converter** (`service/converter.go`) - Pure currency conversion logic
4. **Rate Fetcher** (`service/rate_fetcher.go`) - Main orchestrator with validation
5. **Handler** (`handler/convert_handler.go`) - HTTP request/response handling

## Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)
- API Key from [https://exchangerate.host](https://exchangerate.host/)

## Quick Start

### Option 1: Run with Docker (Recommended)

```bash
# 1. Clone the repository
git clone <your-repo-url>
cd exchange-rate-service

# 2. Set your API key
export API_KEY=your_api_key_here

# 3. Run with Docker Compose
docker-compose up

# Service will be available at http://localhost:8080
```

### Option 2: Run Locally

```bash
# 1. Clone the repository
git clone <your-repo-url>
cd exchange-rate-service

# 2. Install dependencies
go mod download

# 3. Create .env file
echo "API_KEY=your_api_key_here" > .env
echo "PORT=8080" >> .env

# 4. Run the service
go run main.go

# Service will be available at http://localhost:8080
```

## API Documentation

### Convert Currency

Convert an amount from one currency to another.

**Endpoint:** `GET /convert`

**Query Parameters:**

- `from` (required) - Source currency code (USD, INR, EUR, JPY, GBP)
- `to` (required) - Target currency code (USD, INR, EUR, JPY, GBP)
- `amount` (required) - Amount to convert (positive number)
- `date` (optional) - Historical date in YYYY-MM-DD format (max 90 days ago)

**Examples:**

```bash
# Current exchange rate
curl "http://localhost:8080/convert?from=USD&to=INR&amount=100"

# Response:
{
  "from": "USD",
  "to": "INR",
  "amount": 100,
  "result": 8312.50
}

# Historical exchange rate (30 days ago)
curl "http://localhost:8080/convert?from=EUR&to=GBP&amount=50&date=2025-10-28"

# Response:
{
  "from": "EUR",
  "to": "GBP",
  "amount": 50,
  "date": "2025-10-28",
  "result": 43.95
}

# Same currency conversion
curl "http://localhost:8080/convert?from=USD&to=USD&amount=100"

# Response:
{
  "from": "USD",
  "to": "USD",
  "amount": 100,
  "result": 100
}
```

**Error Responses:**

```bash
# Unsupported currency
curl "http://localhost:8080/convert?from=BTC&to=USD&amount=1"
# {"error":"unsupported currency: BTC"}

# Invalid amount
curl "http://localhost:8080/convert?from=USD&to=INR&amount=-100"
# {"error":"amount must be positive"}

# Date too old (>90 days)
curl "http://localhost:8080/convert?from=USD&to=INR&amount=100&date=2024-01-01"
# {"error":"date is too old, maximum lookback is 90 days"}

# Future date
curl "http://localhost:8080/convert?from=USD&to=INR&amount=100&date=2026-01-01"
# {"error":"date cannot be in the future"}

# Missing parameters
curl "http://localhost:8080/convert?from=USD&to=INR"
# {"error":"missing required parameter: amount"}
```

## Testing

### Run All Tests

```bash
# Set API key
export API_KEY=your_api_key_here

# Run tests
go test -v

# Run with coverage
go test -cover
```

### Test Output

```
✓ Basic conversion
✓ All 20 currency pairs
✓ Historical conversion
✓ Unsupported currencies rejected
✓ Invalid amounts rejected
✓ Date validation
✓ Caching performance
```

## Configuration

Environment variables:

```bash
API_KEY=your_api_key_here    # Required: exchangerate-api.com API key
PORT=8080                     # Optional: Server port (default: 8080)
```

## Supported Currencies

- **USD** - United States Dollar
- **INR** - Indian Rupee
- **EUR** - Euro
- **JPY** - Japanese Yen
- **GBP** - British Pound Sterling

## Assumptions

1. **Base Currency:** All exchange rates are relative to USD from the API
2. **Rate Refresh:** Latest rates are refreshed every hour automatically
3. **Historical Data:** Limited to last 90 days as per requirements
4. **Caching Strategy:**
   - Latest rates cached until next hourly refresh
   - Historical rates cached permanently (within 90-day window)
5. **Date Format:** Historical dates must be in YYYY-MM-DD format
6. **Timezone:** All dates are processed in UTC
7. **API Rate Limits:** Free tier API limits are respected
8. **Error Handling:** External API failures return graceful error messages
9. **Concurrency:** Service handles multiple concurrent requests safely using mutex locks
10. **Decimal Precision:** Results are returned as floating-point numbers

## Project Structure

```
exchange-rate-service/
├── main.go                    # Application entry point
├── service/
│   ├── api_client.go         # External API integration
│   ├── cache.go              # In-memory caching
│   ├── converter.go          # Conversion logic
│   └── rate_fetcher.go       # Main service orchestrator
├── handler/
│   └── convert_handler.go    # HTTP handlers
├── service_test.go           # Integration tests
├── Dockerfile                # Docker configuration
├── docker-compose.yml        # Docker Compose configuration
├── .env                      # Environment variables
├── go.mod                    # Go module definition
└── README.md                 # This file
```

## Troubleshooting
1. Make sure that the .env file is located within the /exchange-rate-service
