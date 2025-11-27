package handler

import (
	"net/http"
	"strconv"
	"time"

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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing ~from~ currency",
		})
		return
	}

	if to == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing ~to~ currency",
		})
		return
	}

	if amountStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing ~amount~",
		})
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid amount format, should be a number",
		})
	}

	var date *time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid date format, use YYYY-MM-DD",
			})
			return
		}
		date = &parsedDate
	}

	result, err := h.rateFetcher.ConvertCurrency(from, to, amount, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := gin.H{

		"amount": result,
	}

	c.JSON(http.StatusOK, response)

}
