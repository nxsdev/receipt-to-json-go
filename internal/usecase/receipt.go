package usecase

import (
	"receipt-to-json-go/internal/service"

	"github.com/cockroachdb/errors"
)

type IReceiptUseCase interface {
	ReceiptToJson(imgURL string) (map[string]interface{}, error)
}

type ReceiptUseCase struct {
	ocr service.IOCRService
	cs  service.IClassificationService
}

// NewReceiptUseCase creates a new receipt processor.
func NewReceiptUseCase(ocr service.IOCRService, cs service.IClassificationService) IReceiptUseCase {
	return &ReceiptUseCase{
		ocr: ocr,
		cs:  cs,
	}
}

// ReceiptToJson processes the receipt from the given image URL.
func (p *ReceiptUseCase) ReceiptToJson(imgURL string) (map[string]interface{}, error) {
	ocrResult, err := p.ocr.PerformOCR(imgURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform OCR on the image")
	}

	content, ok := ocrResult["content"].(string)
	if !ok {
		return nil, errors.New("unexpected OCR result format")
	}

	jsonResult, err := p.cs.PerformClassification(content)
	if err != nil {
		return nil, errors.Wrap(err, "failed to classify the OCR content")
	}

	return jsonResult, nil
}
