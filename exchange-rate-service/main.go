package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yourusername/exchange-rate-service/handler"
	"github.com/yourusername/exchange-rate-service/service"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Error .env file not found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rateFetcher := service.NewRateFetcherService()

	rateFetcher.StartHourlyRefresh()

	convertHandler := handler.NewConvertHandler(rateFetcher)
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.GET("/convert", convertHandler.HandleConvert)

	log.Println("Exchange Rate Service Started")
	log.Printf("Server running on port: %s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
