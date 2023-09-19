package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cockroachdb/errors"
)

// IOCRService provides an interface for Optical Character Recognition operations.
type IOCRService interface {
	PerformOCR(imgUrl string) (map[string]interface{}, error)
}

// OCRService is a struct representing the Azure OCR service.
type OCRService struct {
	Client        *http.Client
	VisionURL     string
	VisionHeaders http.Header
	VisionParams  map[string]string
}

// NewOCRService initializes an Azure OCR service with the necessary settings.
func NewOCRService() IOCRService {
	visionURL := os.Getenv("AZURE_VISION_ENDPOINT") + os.Getenv("AZURE_VISION_API_ENDPOINT")
	visionKey := os.Getenv("AZURE_VISION_KEY")
	validateAzureConfig(visionURL, visionKey)

	return &OCRService{
		Client:        &http.Client{Timeout: time.Second * 30},
		VisionURL:     visionURL,
		VisionHeaders: createDefaultHeaders(visionKey),
		VisionParams:  createDefaultParams(),
	}
}

func validateAzureConfig(url, key string) {
	if url == "" || key == "" {
		log.Fatal("Azure configuration is missing. Ensure AZURE_VISION_ENDPOINT, AZURE_VISION_API_ENDPOINT, and AZURE_VISION_KEY are set.")
	}
}

func createDefaultHeaders(key string) http.Header {
	return http.Header{
		"Content-Type":              []string{"application/json"},
		"Ocp-Apim-Subscription-Key": []string{key},
	}
}

func createDefaultParams() map[string]string {
	return map[string]string{
		"features": "read",
		"language": "ja",
	}
}

// PerformOCR initiates the OCR process on the given image URL.
func (s *OCRService) PerformOCR(imgUrl string) (map[string]interface{}, error) {
	req, err := s.createRequest(imgUrl)
	if err != nil {
		return nil, err
	}

	s.setHeadersAndParams(req)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Azure OCR API returned status: %v", resp.Status)
	}

	return parseOCRResponse(resp.Body)
}

func (s *OCRService) createRequest(imgUrl string) (*http.Request, error) {
	data, err := json.Marshal(map[string]string{"url": imgUrl})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal image URL")
	}

	req, err := http.NewRequest(http.MethodPost, s.VisionURL, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	return req, nil
}

func (s *OCRService) setHeadersAndParams(req *http.Request) {
	for key, values := range s.VisionHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	q := req.URL.Query()
	for key, value := range s.VisionParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
}

func parseOCRResponse(body io.Reader) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, errors.Wrap(err, "failed to decode OCR response")
	}

	content, ok := result["readResult"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid OCR response format")
	}

	return content, nil
}
