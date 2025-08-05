package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"prompt2video/internal/models"
	"time"
)

type VideoService struct {
	storageDir string
}

func NewVideoService() *VideoService {
	storageDir := "./storage/videos"
	// Create storage directory if it doesn't exist
	os.MkdirAll(storageDir, 0755)
	
	return &VideoService{
		storageDir: storageDir,
	}
}

func (s *VideoService) RenderVideo(req models.VideoRenderRequest) (*models.VideoRenderResponse, error) {
	// Create temporary directory for this job
	tempDir := filepath.Join(s.storageDir, "temp", req.JobID.String())
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Clean up temp directory after processing
	defer os.RemoveAll(tempDir)

	// Write Manim script to file
	scriptPath := filepath.Join(tempDir, "scene.py")
	if err := s.writeManimScript(scriptPath, req.ManimCode); err != nil {
		return nil, fmt.Errorf("failed to write Manim script: %w", err)
	}

	// Generate output filename
	outputFilename := fmt.Sprintf("video_%s_%d.mp4", req.JobID.String(), time.Now().Unix())
	outputPath := filepath.Join(s.storageDir, outputFilename)

	// Run Manim to generate video
	if err := s.runManim(scriptPath, outputPath, tempDir); err != nil {
		return nil, fmt.Errorf("failed to render video: %w", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("video file was not generated")
	}

	// Return relative URL for serving
	videoURL := fmt.Sprintf("/storage/%s", outputFilename)

	return &models.VideoRenderResponse{
		Success:  true,
		VideoURL: videoURL,
	}, nil
}

func (s *VideoService) writeManimScript(scriptPath, manimCode string) error {
	// Wrap the user's Manim code in a complete script
	fullScript := fmt.Sprintf(`#!/usr/bin/env python3
from manim import *

%s

# Auto-configure scene if no scene class is found
if __name__ == "__main__":
    # This is a fallback if the user didn't create a proper scene
    class GeneratedScene(Scene):
        def construct(self):
            %s
`, manimCode, s.extractConstructMethod(manimCode))

	return os.WriteFile(scriptPath, []byte(fullScript), 0644)
}

func (s *VideoService) extractConstructMethod(manimCode string) string {
	// If the code already contains a Scene class, return empty
	if contains(manimCode, "class") && contains(manimCode, "Scene") {
		return "pass  # Scene already defined above"
	}

	// Otherwise, wrap the code in basic scene construction
	return fmt.Sprintf(`
            # User's code wrapped in scene
            try:
                %s
            except Exception as e:
                # Fallback: show error as text
                error_text = Text(f"Error: {str(e)[:50]}...")
                self.add(error_text)
                self.wait(2)
`, indentCode(manimCode, 12))
}

func (s *VideoService) runManim(scriptPath, outputPath, workDir string) error {
	// Use Docker to run Manim in an isolated environment
	if s.isDockerAvailable() {
		return s.runManimDocker(scriptPath, outputPath, workDir)
	}

	// Fallback: run Manim directly (requires Manim to be installed)
	return s.runManimDirect(scriptPath, outputPath, workDir)
}

func (s *VideoService) isDockerAvailable() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func (s *VideoService) runManimDocker(scriptPath, outputPath, workDir string) error {
	// Create a Docker command to run Manim
	cmd := exec.Command(
		"docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/manim", workDir),
		"-v", fmt.Sprintf("%s:/output", filepath.Dir(outputPath)),
		"manimcommunity/manim:latest",
		"manim", "/manim/scene.py", "-ql", "--output_file", "/output/"+filepath.Base(outputPath),
	)

	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Docker Manim execution failed: %w, output: %s", err, output)
	}

	return nil
}

func (s *VideoService) runManimDirect(scriptPath, outputPath, workDir string) error {
	// Check if manim is available
	if _, err := exec.LookPath("manim"); err != nil {
		return fmt.Errorf("Manim is not installed and Docker is not available")
	}

	// Run manim command directly
	cmd := exec.Command(
		"manim", scriptPath,
		"-ql", // Low quality for faster rendering
		"--output_file", outputPath,
	)

	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Manim execution failed: %w, output: %s", err, output)
	}

	return nil
}

// Helper functions
func contains(text, substr string) bool {
	return len(text) >= len(substr) && findInString(text, substr) != -1
}

func findInString(text, substr string) int {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func indentCode(code string, spaces int) string {
	lines := splitLines(code)
	indent := ""
	for i := 0; i < spaces; i++ {
		indent += " "
	}

	var result string
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		if line != "" {
			result += indent + line
		}
	}
	return result
}

func splitLines(text string) []string {
	var lines []string
	var currentLine string

	for _, char := range text {
		if char == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(char)
		}
	}
	
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	
	return lines
}
