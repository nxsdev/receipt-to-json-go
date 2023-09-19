package service

import (
	"context"
	"encoding/json"
	"receipt-to-json-go/config"
	"strings"

	"github.com/cockroachdb/errors"
	openai "github.com/sashabaranov/go-openai"
)

// IClassificationService is the interface for a service that classifies content.
type IClassificationService interface {
	PerformClassification(content string) (map[string]interface{}, error)
}

// ClassificationService is the default implementation of the Classifier interface using OpenAI.
type ClassificationService struct {
	client *openai.Client
}

// NewDefaultClassifier creates a new default classifier.
func NewClassificationService(client *openai.Client) *ClassificationService {
	return &ClassificationService{client: client}
}

// PerformClassification classifies the given content.
func (c *ClassificationService) PerformClassification(content string) (map[string]interface{}, error) {
	messages := generatePrompt(content)

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "ChatCompletion error")
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("No response from AI")
	}

	return convertToJSON(resp.Choices[0].Message.Content)
}

func generatePrompt(content string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: config.PROMPT_SYSTEM,
		},
		{
			Role:    openai.ChatMessageRoleAssistant,
			Content: config.PROMPT_ASSISTANT + content,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: config.PROMPT_USER + config.PROMPT_EXAMPLE,
		},
	}
}

func convertToJSON(dataString string) (map[string]interface{}, error) {
	startIndex := strings.Index(dataString, "{")
	if startIndex == -1 {
		return nil, errors.New("could not find start of JSON in string")
	}

	endIndex := strings.LastIndex(dataString, "}")
	if endIndex == -1 {
		return nil, errors.New("could not find end of JSON in string")
	}

	validJSONString := dataString[startIndex : endIndex+1]

	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(validJSONString), &jsonData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal JSON data")
	}

	return jsonData, nil
}
