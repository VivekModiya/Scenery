package services

import (
	"testing"
)

func TestLLMServiceNew(t *testing.T) {
	service := NewLLMService("test-openai-key", "test-gemini-key")
	
	if service == nil {
		t.Error("Expected LLMService to be created, got nil")
	}
	
	if service.geminiAPIKey != "test-gemini-key" {
		t.Errorf("Expected geminiAPIKey to be 'test-gemini-key', got %s", service.geminiAPIKey)
	}
}

func TestBuildSystemPrompt(t *testing.T) {
	service := NewLLMService("", "")
	
	prompt := service.buildSystemPrompt("mathematics")
	
	if prompt == "" {
		t.Error("Expected system prompt to be generated, got empty string")
	}
	
	// Check if subject is included in prompt
	if !contains(prompt, "mathematics") {
		t.Error("Expected system prompt to contain subject 'mathematics'")
	}
	
	// Check for required sections
	if !contains(prompt, "EXPLANATION:") {
		t.Error("Expected system prompt to contain 'EXPLANATION:' marker")
	}
	
	if !contains(prompt, "MANIM_CODE:") {
		t.Error("Expected system prompt to contain 'MANIM_CODE:' marker")
	}
}

func TestBuildUserPrompt(t *testing.T) {
	service := NewLLMService("", "")
	
	userPrompt := service.buildUserPrompt("Explain photosynthesis")
	
	if userPrompt == "" {
		t.Error("Expected user prompt to be generated, got empty string")
	}
	
	if !contains(userPrompt, "photosynthesis") {
		t.Error("Expected user prompt to contain 'photosynthesis'")
	}
}

func TestParseResponse(t *testing.T) {
	service := NewLLMService("", "")
	
	testResponse := `EXPLANATION:
This is a test explanation of a concept.

MANIM_CODE:
` + "```python" + `
from manim import *

class TestScene(Scene):
    def construct(self):
        text = Text("Hello World")
        self.add(text)
        self.wait(2)
` + "```"
	
	explanation, manimCode, err := service.parseResponse(testResponse)
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if explanation == "" {
		t.Error("Expected explanation to be extracted, got empty string")
	}
	
	if manimCode == "" {
		t.Error("Expected Manim code to be extracted, got empty string")
	}
	
	if !contains(explanation, "test explanation") {
		t.Error("Expected explanation to contain 'test explanation'")
	}
	
	if !contains(manimCode, "TestScene") {
		t.Error("Expected Manim code to contain 'TestScene'")
	}
}

func TestParseResponseInvalidFormat(t *testing.T) {
	service := NewLLMService("", "")
	
	invalidResponse := "This is not a properly formatted response"
	
	_, _, err := service.parseResponse(invalidResponse)
	
	if err == nil {
		t.Error("Expected error for invalid response format, got nil")
	}
}
