package router

import (
	"receipt-to-json-go/internal/delivery/http/controller"

	"github.com/labstack/echo/v4"
)

// SetupReceiptRoutes sets up the routes for the receipt controller.
func SetupReceiptRoutes(e *echo.Echo, rc *controller.ReceiptController) {
	e.POST("/api/receipt/process", rc.ProcessReceipt)
}
