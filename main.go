package main

import (
	"fmt"
	"os"
	_ "receipt-to-json-go/docs"
	"receipt-to-json-go/internal/delivery/http/controller"
	"receipt-to-json-go/internal/delivery/http/router"
	"receipt-to-json-go/internal/service"
	"receipt-to-json-go/internal/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
	echoswagger "github.com/swaggo/echo-swagger"
)

// @title Receipt to JSON API
// @discription This is an API that converts receipt information from an image URL to JSON.
// @version 0.1
// @host localhost:8080
func main() {
	e := echo.New()
	//

	e.GET("/swagger/*", echoswagger.WrapHandler)

	ocrService := service.NewOCRService()
	classificationService := service.NewClassificationService(openai.NewClient(os.Getenv("OPENAI_KEY")))
	receiptProcessor := usecase.NewReceiptUseCase(ocrService, classificationService)
	receiptController := controller.NewReceiptController(receiptProcessor)
	router.SetupReceiptRoutes(e, receiptController)

	e.Start(":8080")
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("No .env file found: %v", err)
	}
}
