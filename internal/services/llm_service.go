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
	return fmt.Sprintf(`You are an expert educator and animator who creates educational content using Manim (Mathematical Animation Engine).

Your task is to:
1. Provide a clear, concise explanation of the concept
2. Generate Manim Python code that creates a 15-second educational animation

Guidelines:
- Keep explanations under 200 words
- Focus on %s concepts when relevant
- Create visually engaging animations
- Use appropriate Manim classes (Text, Circle, Square, Arrow, etc.)
- Include smooth transitions and animations
- Ensure the animation runs for approximately 15 seconds
- Use clear, readable fonts and colors

Return your response in this exact format:
EXPLANATION:
[Your explanation here]

MANIM_CODE:
` + "```python" + `
[Your Manim code here]
` + "```" + ``, subject)
}

func (s *LLMService) buildUserPrompt(prompt string) string {
	return fmt.Sprintf(`Create an educational animation explaining: %s

Make sure the animation is engaging, informative, and suitable for learning. The animation should be approximately 15 seconds long.`, prompt)
}

func (s *LLMService) parseResponse(content string) (explanation, manimCode string, err error) {
	// Simple parsing - look for EXPLANATION: and MANIM_CODE: markers
	explanationStart := findMarker(content, "EXPLANATION:")
	manimStart := findMarker(content, "MANIM_CODE:")

	if explanationStart == -1 || manimStart == -1 {
		return "", "", fmt.Errorf("response format invalid - missing markers")
	}

	// Extract explanation
	explanationEnd := manimStart
	if explanationStart+12 < len(content) && explanationEnd > explanationStart+12 {
		explanation = content[explanationStart+12 : explanationEnd]
		explanation = trimWhitespace(explanation)
	}

	// Extract Manim code - look for the code block after MANIM_CODE:
	remainingContent := content[manimStart:]
	codeStart := findMarker(remainingContent, "```python")
	if codeStart == -1 {
		return "", "", fmt.Errorf("could not find python code block")
	}
	
	// Find the closing ```
	searchStart := codeStart + 9 // Skip "```python"
	codeEnd := findMarker(remainingContent[searchStart:], "```")
	if codeEnd == -1 {
		return "", "", fmt.Errorf("could not find end of code block")
	}

	codeStartIdx := codeStart + 9 // Skip "```python"
	codeEndIdx := searchStart + codeEnd
	if codeStartIdx < len(remainingContent) && codeEndIdx <= len(remainingContent) {
		manimCode = remainingContent[codeStartIdx:codeEndIdx]
		manimCode = trimWhitespace(manimCode)
	}

	if explanation == "" || manimCode == "" {
		return "", "", fmt.Errorf("failed to extract explanation or Manim code")
	}

	return explanation, manimCode, nil
}

func findMarker(text, marker string) int {
	for i := 0; i <= len(text)-len(marker); i++ {
		if text[i:i+len(marker)] == marker {
			return i
		}
	}
	return -1
}

func trimWhitespace(text string) string {
	start := 0
	end := len(text)

	// Trim leading whitespace
	for start < len(text) && (text[start] == ' ' || text[start] == '\n' || text[start] == '\t' || text[start] == '\r') {
		start++
	}

	// Trim trailing whitespace
	for end > start && (text[end-1] == ' ' || text[end-1] == '\n' || text[end-1] == '\t' || text[end-1] == '\r') {
		end--
	}

	return text[start:end]
}
