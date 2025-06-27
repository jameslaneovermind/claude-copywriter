package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// processFile processes a single file through Claude API
func processFile(inputFile, outputFile, role, systemPrompt, customPrompt, sharedContent string, client *ClaudeClient, verbose bool) (*ProcessingResult, error) {
	result := &ProcessingResult{
		InputFile:  inputFile,
		OutputFile: outputFile,
		Role:       role,
	}

	if verbose {
		fmt.Printf("📖 Reading file: %s\n", inputFile)
	}

	// Read input file
	content, err := readFileContent(inputFile)
	if err != nil {
		result.Error = fmt.Errorf("failed to read input file: %w", err)
		return result, result.Error
	}

	// Combine with shared content and custom prompt if provided
	finalContent := combineContent(content, sharedContent, customPrompt, verbose)

	if verbose {
		fmt.Printf("📝 Processing with Claude (role: %s)...\n", role)
	}

	// Process through Claude API
	response, err := client.ProcessContent(finalContent, systemPrompt)
	if err != nil {
		result.Error = fmt.Errorf("Claude API processing failed: %w", err)
		return result, result.Error
	}

	// Extract processed content
	processedContent := client.ExtractTextFromResponse(response)
	result.TokensUsed = client.GetTokenUsage(response)

	if verbose {
		fmt.Printf("💾 Writing output: %s\n", outputFile)
	}

	// Write output file
	if err := writeFileContent(outputFile, processedContent); err != nil {
		result.Error = fmt.Errorf("failed to write output file: %w", err)
		return result, result.Error
	}

	result.Success = true
	return result, nil
}

// processBatch processes multiple files in a directory
func processBatch(batchDir, role, systemPrompt, customPrompt, sharedContent string, client *ClaudeClient, maxConcurrent int, verbose bool) ([]ProcessingResult, error) {
	// Find all text files in directory
	files, err := findTextFiles(batchDir)
	if err != nil {
		return nil, fmt.Errorf("failed to find files: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no text files found in directory: %s", batchDir)
	}

	if verbose {
		fmt.Printf("🔍 Found %d files to process\n", len(files))
	}

	// Create channels for concurrent processing
	fileChan := make(chan string, len(files))
	resultChan := make(chan ProcessingResult, len(files))

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for inputFile := range fileChan {
				outputFile := generateOutputFilename(inputFile, role)
				
				// Add small delay to avoid hitting rate limits
				time.Sleep(100 * time.Millisecond)
				
				result, _ := processFile(inputFile, outputFile, role, systemPrompt, customPrompt, sharedContent, client, verbose)
				resultChan <- *result
			}
		}()
	}

	// Send files to workers
	go func() {
		for _, file := range files {
			fileChan <- file
		}
		close(fileChan)
	}()

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var results []ProcessingResult
	for result := range resultChan {
		results = append(results, result)
	}

	return results, nil
}

// readFileContent reads content from a file, handling different file types
func readFileContent(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Convert to string and validate it's text
	text := string(content)
	if !isTextContent(text) {
		return "", fmt.Errorf("file appears to contain binary data: %s", filename)
	}

	return text, nil
}

// writeFileContent writes content to a file
func writeFileContent(filename, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// findTextFiles finds all text files in a directory
func findTextFiles(dir string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if it's a text file based on extension
		if isTextFile(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// isTextFile checks if a file is likely a text file based on extension
func isTextFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	textExtensions := map[string]bool{
		".txt":  true,
		".md":   true,
		".text": true,
		".csv":  true,
		".json": true,
		".xml":  true,
		".html": true,
		".htm":  true,
		".css":  true,
		".js":   true,
		".py":   true,
		".go":   true,
		".java": true,
		".cpp":  true,
		".c":    true,
		".h":    true,
		".sql":  true,
		".yaml": true,
		".yml":  true,
		".toml": true,
		".ini":  true,
		".cfg":  true,
		".conf": true,
		".log":  true,
	}

	return textExtensions[ext]
}

// isTextContent performs a basic check to see if content is text
func isTextContent(content string) bool {
	// Simple heuristic: if more than 90% of characters are printable or valid Unicode, consider it text
	if len(content) == 0 {
		return true
	}

	// Convert to runes to properly handle Unicode
	runes := []rune(content)
	validCount := 0
	
	for _, r := range runes {
		// Accept ASCII printable characters, whitespace, and most Unicode characters
		if (r >= 32 && r <= 126) || // ASCII printable
		   r == '\n' || r == '\r' || r == '\t' || // Common whitespace
		   (r > 126 && r < 0xFFFE && r != 0xFFFD) { // Valid Unicode (excluding replacement char)
			validCount++
		}
	}

	ratio := float64(validCount) / float64(len(runes))
	return ratio > 0.90 // Lowered threshold to be more permissive
}

// combineContent combines the main content with shared content and custom prompt
func combineContent(mainContent, sharedContent, customPrompt string, verbose bool) string {
	var sections []string

	// Add custom prompt instructions if provided
	if customPrompt != "" {
		if verbose {
			fmt.Printf("📝 Adding custom prompt instructions\n")
		}
		sections = append(sections, fmt.Sprintf("SPECIFIC INSTRUCTIONS FOR THIS REQUEST:\n%s", customPrompt))
	}

	// Add shared content if provided
	if sharedContent != "" {
		if verbose {
			fmt.Printf("🔗 Including shared context information\n")
		}
		sections = append(sections, fmt.Sprintf("SHARED CONTEXT INFORMATION:\n%s", sharedContent))
	}

	// Add main content
	sections = append(sections, fmt.Sprintf("MAIN CONTENT TO PROCESS:\n%s", mainContent))

	// Combine all sections with clear separators
	return strings.Join(sections, "\n\n---\n\n")
}