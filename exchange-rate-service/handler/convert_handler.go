package handler

import (
	"net/http"
	"strconv"
	"time"

	appErrors "github.com/yourusername/exchange-rate-service/errors"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/exchange-rate-service/service"
)

type ConvertHandler struct {
	rateFetcher *service.RateFetcherService
}

func NewConvertHandler(rateFetcher *service.RateFetcherService) *ConvertHandler {
	return &ConvertHandler{
		rateFetcher: rateFetcher,
	}
}

func (h *ConvertHandler) HandleConvert(c *gin.Context) {

	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")
	dateStr := c.Query("date")

	if from == "" {
		h.respondWithError(c, appErrors.MissingParameterError("from"))
		return
	}

	if to == "" {
		h.respondWithError(c, appErrors.MissingParameterError("to"))
		return
	}

	if amountStr == "" {
		h.respondWithError(c, appErrors.MissingParameterError("amount"))
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		h.respondWithError(c, appErrors.InvalidAmountError())
		return
	}

	var date *time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			h.respondWithError(c, appErrors.InvalidDateFormatError())
			return
		}
		date = &parsedDate
	}

	result, err := h.rateFetcher.ConvertCurrency(from, to, amount, date)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	response := gin.H{

		"amount": result,
	}

	c.JSON(http.StatusOK, response)

}

func (h *ConvertHandler) respondWithError(c *gin.Context, err error) {

	customErr, ok := err.(*appErrors.CustomError)

	if ok {
		c.JSON(customErr.GetHTTPStatus(), gin.H{
			"error":        customErr.Code,
			"errorMessage": customErr.Message,
		})

		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "internal server error",
	})

}
