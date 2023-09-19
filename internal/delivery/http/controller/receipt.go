package controller

import (
	"net/http"
	"receipt-to-json-go/internal/usecase"

	"github.com/labstack/echo/v4"
)

type ReceiptController struct {
	uc usecase.IReceiptUseCase
}

// NewReceiptController creates a new receipt controller.
func NewReceiptController(rp usecase.IReceiptUseCase) *ReceiptController {
	return &ReceiptController{uc: rp}
}

// @Summary Process receipt from image URL
// @Description Converts receipt information from an image URL to JSON
// @Accept json
// @Produce json
// @Param image_url query string true "URL of the image containing the receipt"
// @Success 200 {object} map[string]interface{}
// @Router /api/receipt/process [post]
func (rc *ReceiptController) ProcessReceipt(c echo.Context) error {
	var input struct {
		ImageURL string `json:"image_url"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result, err := rc.uc.ReceiptToJson(input.ImageURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
