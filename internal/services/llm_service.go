package services

import (
	"context"
	"fmt"
	"prompt2video/internal/models"

	"github.com/sashabaranov/go-openai"
)

type LLMService struct {
	openaiClient *openai.Client
	geminiAPIKey string
}

func NewLLMService(openaiAPIKey, geminiAPIKey string) *LLMService {
	var openaiClient *openai.Client
	if openaiAPIKey != "" {
		openaiClient = openai.NewClient(openaiAPIKey)
	}

	return &LLMService{
		openaiClient: openaiClient,
		geminiAPIKey: geminiAPIKey,
	}
}

func (s *LLMService) GenerateExplanation(prompt, subject string) (*models.LLMResponse, error) {
	if s.openaiClient == nil {
		return nil, fmt.Errorf("no LLM client configured")
	}

	// Create the system prompt for generating Manim code
	systemPrompt := s.buildSystemPrompt(subject)
	userPrompt := s.buildUserPrompt(prompt)

	resp, err := s.openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			MaxTokens:   2000,
			Temperature: 0.7,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	// Parse the response to extract explanation and Manim code
	explanation, manimCode, err := s.parseResponse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	return &models.LLMResponse{
		Explanation: explanation,
		ManimCode:   manimCode,
		Success:     true,
	}, nil
}

func (s *LLMService) buildSystemPrompt(subject string) string {
	// Create a system prompt that instructs the LLM to generate educational animations
}

func (s *LLMService) buildUserPrompt(prompt string) string {
// Create a user prompt that includes the subject and the user's request
}

func (s *LLMService) parseResponse(content string) (explanation, manimCode string, err error) {
// Parse the LLM response to extract the explanation and Manim code
}


