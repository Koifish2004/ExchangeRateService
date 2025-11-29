# Exchange Rate Service

A currency exchange rate service built in Go that provides real-time and historical currency conversion with in-memory caching.

## Features

- Real-time currency conversion for USD, INR, EUR, JPY, GBP, BTC
- Historical exchange rates (up to 90 days)
- In-memory caching for performance
- Hourly automatic rate refresh
- Thread-safe concurrent request handling
- RESTful API with comprehensive validation

## Prerequisites

- Go 1.21 or higher
- API Key from [exchangerate.host](https://exchangerate.host/)
- Docker (optional)

## Quick Start

### Local Setup

```bash
# Clone repository
git clone <your-repo-url>
cd exchange-rate-service

# Install dependencies
go mod download

# Set environment variables
export API_KEY=your_api_key_here
export PORT=8080

# Run service
go run main.go
```

### Docker Setup

```bash
# Set API key
export API_KEY=your_api_key_here

# Run with Docker Compose
docker-compose up
```

Service runs at `http://localhost:8080`

## API Documentation

### Convert Currency

**Endpoint:** `GET /convert`

**Query Parameters:**

| Parameter | Required | Description                                    |
| --------- | -------- | ---------------------------------------------- |
| `from`    | Yes      | Source currency (USD, INR, EUR, JPY, GBP, BTC) |
| `to`      | Yes      | Target currency (USD, INR, EUR, JPY, GBP, BTC) |
| `amount`  | Yes      | Amount to convert (positive number)            |
| `date`    | No       | Historical date (YYYY-MM-DD, max 90 days ago)  |

**Success Response:**

```json
{
  "amount": "100"
}
```

**Example Requests:**

```bash
# Current rate
curl "http://localhost:8080/convert?from=USD&to=INR&amount=100"

# Historical rate
curl "http://localhost:8080/convert?from=EUR&to=GBP&amount=50&date=2025-10-28"

# Same currency
curl "http://localhost:8080/convert?from=USD&to=USD&amount=100"
```

**Error Response:**

```json
{
  "error": "UNSUPPORTED_CURRENCY",
  "errorMessage": "unsupported currency: XYZ"
}
```

## Configuration

Environment variables:

```bash
API_KEY=your_api_key_here    # Required
PORT=8080                     # Optional (default: 8080)
```

## Architecture

The service consists of five main components:

1. **API Client** - Fetches rates from exchangerate.host
2. **Cache** - Thread-safe in-memory storage with mutex locks
3. **Converter** - Currency conversion calculations using decimal precision
4. **Rate Fetcher** - Orchestrates validation, caching, and conversion
5. **Handler** - HTTP request/response processing with Gin framework

## Supported Currencies

- USD - United States Dollar
- INR - Indian Rupee
- EUR - Euro
- JPY - Japanese Yen
- GBP - British Pound Sterling
- BTC - Bitcoin

## Project Structure

```
exchange-rate-service/
├── main.go                    # Application entry point
├── handler/
│   └── convert_handler.go    # HTTP request handlers
├── service/
│   ├── api_client.go         # External API integration
│   ├── cache.go              # In-memory caching
│   ├── converter.go          # Conversion logic
│   └── rate_fetcher.go       # Service orchestrator
├── errors/
│   └── errors.go             # Custom error types
├── Dockerfile                # Container configuration
├── docker-compose.yml        # Docker Compose setup
└── go.mod                    # Go module dependencies
```

## Testing

```bash
# Set API key
export API_KEY=your_api_key_here

# Run tests
go test ./... -v

# Run with coverage
go test ./... -cover
```

## Dependencies

- **gin-gonic/gin** - HTTP web framework
- **joho/godotenv** - Environment variable management
- **shopspring/decimal** - Precise decimal arithmetic
