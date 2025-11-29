package errors

import "fmt"

type ErrorCode string

const (
	ErrMissingParameter    ErrorCode = "MISSING_PARAMETER"
	ErrInvalidAmount       ErrorCode = "INVALID_AMOUNT"
	ErrInvalidDate         ErrorCode = "INVALID_DATE_FORMAT"
	ErrUnsupportedCurrency ErrorCode = "UNSUPPORTED_CURRENCY"
	ErrDateTooOld          ErrorCode = "DATE_TOO_OLD"
	ErrFutureDate          ErrorCode = "FUTURE_DATE"

	ErrAPIFetchFailed ErrorCode = "API_FETCH_FAILED"
	ErrAPIBadStatus   ErrorCode = "API_BAD_STATUS"
	ErrAPIBadResponse ErrorCode = "API_BAD_RESPONSE"

	ErrMissingRate      ErrorCode = "MISSING_EXCHANGE_RATE"
	ErrInvalidRate      ErrorCode = "INVALID_EXCHANGE_RATE"
	ErrConversionFailed ErrorCode = "CONVERSION_FAILED"
)

type ErrorCategory string

const (
	CategoryValidation ErrorCategory = "VALIDATION_ERROR"
	CategoryAPI        ErrorCategory = "API_ERROR"
	CategoryInternal   ErrorCategory = "INTERNAL_ERROR"
)

type CustomError struct {
	Code     ErrorCode
	Category ErrorCategory
	Message  string
	Err      error
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func (c CustomError) ErrorMessage() string {
	return c.Message
}

func (e *CustomError) GetHTTPStatus() int {
	switch e.Category {
	case CategoryValidation:
		return 400
	case CategoryAPI:
		return 502
	case CategoryInternal:
		return 500
	default:
		return 500
	}
}

func newCustomError(code ErrorCode, category ErrorCategory, message string, err error) *CustomError {
	return &CustomError{
		Code:     code,
		Category: category,
		Message:  message,
		Err:      err,
	}
}

func NewAPIError(message string, err error) *CustomError {
	return newCustomError(
		ErrAPIBadResponse,
		CategoryAPI,
		message,
		err,
	)
}

func MissingParameterError(param string) *CustomError {
	return newCustomError(
		ErrMissingParameter,
		CategoryValidation,
		fmt.Sprintf("missing required parameter: %s", param),
		nil,
	)
}

func InvalidAmountError() *CustomError {
	return newCustomError(
		ErrInvalidAmount,
		CategoryValidation,
		"amount must be positive",
		nil,
	)
}

func InvalidDateFormatError() *CustomError {
	return newCustomError(
		ErrInvalidDate,
		CategoryValidation,
		"invalid date format, use YYYY-MM-DD",
		nil,
	)
}

func UnsupportedCurrencyError(currency string) *CustomError {
	return newCustomError(
		ErrUnsupportedCurrency,
		CategoryValidation,
		fmt.Sprintf("unsupported currency: %s", currency),
		nil,
	)
}

func DateTooOldError() *CustomError {
	return newCustomError(
		ErrDateTooOld,
		CategoryValidation,
		"date is too old, maximum lookback is 90 days",
		nil,
	)
}

func FutureDateError() *CustomError {
	return newCustomError(
		ErrFutureDate,
		CategoryValidation,
		"date cannot be in the future",
		nil,
	)
}

//api errors

func APIFetchError(err error) *CustomError {
	return newCustomError(
		ErrAPIFetchFailed,
		CategoryAPI,
		"failed to fetch exchange rates from external API",
		err,
	)
}

func APIBadStatusError(statusCode int) *CustomError {
	return newCustomError(
		ErrAPIBadStatus,
		CategoryAPI,
		fmt.Sprintf("API returned status code: %d", statusCode),
		nil,
	)
}

func APIResponseError(err error) *CustomError {
	return newCustomError(
		ErrAPIBadResponse,
		CategoryAPI,
		"failed to parse API response",
		err,
	)
}

//internal service error

func MissingRateError(currency string) *CustomError {
	return newCustomError(
		ErrMissingRate,
		CategoryInternal,
		fmt.Sprintf("exchange rate not available for: %s", currency),
		nil,
	)
}

func InvalidRateError(currency string) *CustomError {
	return newCustomError(
		ErrInvalidRate,
		CategoryInternal,
		fmt.Sprintf("invalid exchange rate (zero) for: %s", currency),
		nil,
	)
}

func ConversionError(err error) *CustomError {
	return newCustomError(
		ErrConversionFailed,
		CategoryInternal,
		"currency conversion failed",
		err,
	)
}
