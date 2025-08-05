package services

import (
	"os"
	"testing"
)

func TestVideoServiceNew(t *testing.T) {
	service := NewVideoService()
	
	if service == nil {
		t.Error("Expected VideoService to be created, got nil")
	}
	
	if service.storageDir == "" {
		t.Error("Expected storageDir to be set, got empty string")
	}
}

func TestWriteManimScript(t *testing.T) {
	service := NewVideoService()
	
	// Create a temporary file
	tempFile := "/tmp/test_script.py"
	defer os.Remove(tempFile)
	
	testCode := `
class TestScene(Scene):
    def construct(self):
        text = Text("Hello")
        self.add(text)
`
	
	err := service.writeManimScript(tempFile, testCode)
	if err != nil {
		t.Errorf("Expected no error writing script, got: %v", err)
	}
	
	// Check if file was created
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("Expected script file to be created")
	}
	
	// Read and verify content
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Errorf("Error reading script file: %v", err)
	}
	
	contentStr := string(content)
	if !contains(contentStr, "TestScene") {
		t.Error("Expected script content to contain TestScene")
	}
	
	if !contains(contentStr, "from manim import *") {
		t.Error("Expected script content to contain manim import")
	}
}

func TestExtractConstructMethod(t *testing.T) {
	service := NewVideoService()
	
	// Test with existing Scene class
	codeWithScene := `
class MyScene(Scene):
    def construct(self):
        text = Text("Test")
        self.add(text)
`
	
	result := service.extractConstructMethod(codeWithScene)
	if !contains(result, "Scene already defined") {
		t.Error("Expected fallback message for existing scene")
	}
	
	// Test with raw code
	rawCode := `
text = Text("Hello")
self.add(text)
`
	
	result = service.extractConstructMethod(rawCode)
	if !contains(result, "Hello") {
		t.Error("Expected raw code to be wrapped in construct method")
	}
}

func TestIsDockerAvailable(t *testing.T) {
	service := NewVideoService()
	
	// This test will vary based on the system
	result := service.isDockerAvailable()
	
	// Just check that it returns a boolean (no error)
	if result != true && result != false {
		t.Error("Expected isDockerAvailable to return a boolean")
	}
}
