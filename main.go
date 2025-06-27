package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the application configuration
type Config struct {
	APIKey              string            `json:"api_key"`
	BaseURL             string            `json:"base_url"`
	Model               string            `json:"model"`
	PerplexityAPIKey    string            `json:"perplexity_api_key"`
	PerplexityModel     string            `json:"perplexity_model"`
	PerplexityMaxTokens int               `json:"perplexity_max_tokens"`
	Roles               map[string]string `json:"roles"`
}

// ProcessingResult holds the result of processing a file
type ProcessingResult struct {
	InputFile   string
	OutputFile  string
	Role        string
	Success     bool
	Error       error
	TokensUsed  int
}

func main() {
	// Command line flags
	var (
		inputFile     = flag.String("input", "", "Input file path (required)")
		outputFile    = flag.String("output", "", "Output file path (optional, defaults to input_[role].ext)")
		role          = flag.String("role", "copywriter", "Role/prompt to use")
		customPrompt  = flag.String("prompt", "", "Additional custom instructions for this specific request")
		configFile    = flag.String("config", "config.json", "Configuration file path")
		listRoles     = flag.Bool("list-roles", false, "List available roles")
		batchDir      = flag.String("batch", "", "Process all files in directory")
		concurrent    = flag.Int("concurrent", 3, "Number of concurrent API calls")
		verbose       = flag.Bool("verbose", false, "Verbose output")
		includeShared = flag.String("include", "", "Include shared content (filename or comma-separated list)")
		sharedDir     = flag.String("shared-dir", "shared", "Directory containing shared content files")
		listShared    = flag.Bool("list-shared", false, "List available shared content files")
		research      = flag.String("research", "", "Perform research using Perplexity API before processing")
		researchOnly  = flag.Bool("research-only", false, "Only perform research, don't process with Claude")
	)
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Claude client
	client := NewClaudeClient(config.APIKey, config.BaseURL, config.Model)

	// Create Perplexity client if research is needed
	var perplexityClient *PerplexityClient
	if *research != "" || *researchOnly {
		perplexityClient = NewPerplexityClient(config.PerplexityAPIKey, "", config.PerplexityModel, config.PerplexityMaxTokens)
	}

	// Handle list shared content command
	if *listShared {
		listSharedContent(*sharedDir)
		return
	}

	// Handle list roles command
	if *listRoles {
		listAvailableRoles(config.Roles)
		return
	}

	// Validate role exists
	systemPrompt, exists := config.Roles[*role]
	if !exists {
		log.Fatalf("Role '%s' not found. Use -list-roles to see available roles.", *role)
	}

	if *verbose {
		fmt.Printf("Using role: %s\n", *role)
		fmt.Printf("Model: %s\n", config.Model)
		if *includeShared != "" {
			fmt.Printf("Including shared content: %s\n", *includeShared)
		}
		if *customPrompt != "" {
			fmt.Printf("Custom prompt: %s\n", truncateString(*customPrompt, 100))
		}
		if *research != "" {
			fmt.Printf("Research query: %s\n", truncateString(*research, 100))
		}
	}

	// Perform research if requested
	var researchContent string
	if *research != "" && perplexityClient != nil {
		if *verbose {
			fmt.Printf("🔍 Performing research with Perplexity...\n")
		}
		
		researchResponse, err := perplexityClient.ResearchContent(*research, "You are a research assistant. Provide comprehensive, accurate information with proper citations.")
		if err != nil {
			log.Fatalf("Research failed: %v", err)
		}
		
		researchContent = perplexityClient.FormatResearchOutput(researchResponse)
		
		if *verbose {
			fmt.Printf("📊 Research tokens used: %d\n", perplexityClient.GetTokenUsage(researchResponse))
		}
		
		// If research-only mode, save research and exit
		if *researchOnly {
			if *inputFile == "" {
				log.Fatal("Input file is required even in research-only mode to determine output filename")
			}
			
			if *outputFile == "" {
				*outputFile = generateOutputFilename(*inputFile, "research")
			}
			
			if err := writeFileContent(*outputFile, researchContent); err != nil {
				log.Fatalf("Failed to write research output: %v", err)
			}
			
			fmt.Printf("✅ Research completed: %s\n", *outputFile)
			return
		}
	}

	// Load shared content if specified
	var sharedContent string
	if *includeShared != "" {
		var err error
		sharedContent, err = loadSharedContent(*includeShared, *sharedDir, *verbose)
		if err != nil {
			log.Fatalf("Failed to load shared content: %v", err)
		}
	}
	
	// Add research content to shared content if research was performed
	if researchContent != "" {
		if sharedContent != "" {
			sharedContent += "\n\n---\n\n"
		}
		sharedContent += fmt.Sprintf("RESEARCH FINDINGS:\n%s", researchContent)
		if *verbose {
			fmt.Printf("🔬 Including research findings in context\n")
		}
	}

	// Handle batch processing
	if *batchDir != "" {
		results, err := processBatch(*batchDir, *role, systemPrompt, *customPrompt, sharedContent, client, *concurrent, *verbose)
		if err != nil {
			log.Fatalf("Batch processing failed: %v", err)
		}
		printBatchResults(results)
		return
	}

	// Handle single file processing
	if *inputFile == "" {
		flag.Usage()
		log.Fatal("Input file is required")
	}

	// Generate output filename if not provided
	if *outputFile == "" {
		*outputFile = generateOutputFilename(*inputFile, *role)
	}

	// Process single file
	result, err := processFile(*inputFile, *outputFile, *role, systemPrompt, *customPrompt, sharedContent, client, *verbose)
	if err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	if result.Success {
		fmt.Printf("✅ Successfully processed: %s -> %s\n", result.InputFile, result.OutputFile)
		if result.TokensUsed > 0 {
			fmt.Printf("📊 Tokens used: %d\n", result.TokensUsed)
		}
	} else {
		fmt.Printf("❌ Failed to process: %s - %v\n", result.InputFile, result.Error)
	}
}

// loadConfig loads configuration from JSON file
func loadConfig(filename string) (*Config, error) {
	// Create default config if file doesn't exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		defaultConfig := getDefaultConfig()
		return defaultConfig, saveConfig(filename, defaultConfig)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	// Use environment variable if API key not in config
	if config.APIKey == "" {
		config.APIKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	
	// Use environment variable for Perplexity API key if not in config
	if config.PerplexityAPIKey == "" {
		config.PerplexityAPIKey = os.Getenv("PERPLEXITY_API_KEY")
	}

	return &config, nil
}

// saveConfig saves configuration to JSON file
func saveConfig(filename string, config *Config) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

// listAvailableRoles prints all available roles
func listAvailableRoles(roles map[string]string) {
	fmt.Println("Available roles:")
	for role, prompt := range roles {
		fmt.Printf("\n🎭 %s:\n   %s\n", role, truncateString(prompt, 100))
	}
}

// generateOutputFilename creates an output filename based on input and role
func generateOutputFilename(inputFile, role string) string {
	ext := filepath.Ext(inputFile)
	base := strings.TrimSuffix(inputFile, ext)
	return fmt.Sprintf("%s_%s%s", base, role, ext)
}

// truncateString truncates a string to specified length
func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// printBatchResults prints summary of batch processing results
func printBatchResults(results []ProcessingResult) {
	successful := 0
	failed := 0
	totalTokens := 0

	for _, result := range results {
		if result.Success {
			successful++
			totalTokens += result.TokensUsed
			fmt.Printf("✅ %s -> %s\n", result.InputFile, result.OutputFile)
		} else {
			failed++
			fmt.Printf("❌ %s: %v\n", result.InputFile, result.Error)
		}
	}

	fmt.Printf("\n📊 Batch Processing Summary:\n")
	fmt.Printf("   Successful: %d\n", successful)
	fmt.Printf("   Failed: %d\n", failed)
	fmt.Printf("   Total tokens used: %d\n", totalTokens)
}

// getDefaultConfig returns a default configuration
func getDefaultConfig() *Config {
	return &Config{
		APIKey:              "",
		BaseURL:             "https://api.anthropic.com",
		Model:               "claude-3-sonnet-20240229",
		PerplexityAPIKey:    "",
		PerplexityModel:     "llama-3.1-sonar-large-128k-online",
		PerplexityMaxTokens: 8192, // High default for detailed research
		Roles: map[string]string{
			"copywriter": "You are an expert copywriter with 10+ years of experience. Transform the input content into compelling, engaging copy that drives action. Focus on clarity, persuasion, and emotional connection. Use proven copywriting frameworks like AIDA, PAS, or Before-After-Bridge when appropriate. Make the content more engaging while maintaining accuracy.",
			"researcher": "You are a thorough researcher and fact-checker. Analyze the input content and provide detailed research insights, additional context, fact-checking, and citations where appropriate. Identify gaps in information, suggest areas for deeper investigation, and provide relevant background information that adds value to the original content.",
			"reviewer": "You are a critical content reviewer and editor. Evaluate the input content for accuracy, clarity, tone, structure, and overall effectiveness. Provide constructive feedback and specific suggestions for improvement. Check for logical flow, consistency, grammar, and readability. Highlight strengths and areas for enhancement.",
			"formatter": "You are a formatting and structure specialist. Improve the organization, readability, and visual presentation of the content while maintaining its core message. Add appropriate headings, bullet points, spacing, and structure. Optimize for scanning and readability. Ensure consistent formatting throughout.",
			"summarizer": "You are a summarization expert skilled at distilling complex information. Create a concise, comprehensive summary that captures the key points, main ideas, and essential takeaways of the input content. Maintain the original intent while making it more digestible and actionable.",
			"editor": "You are a professional editor with expertise in multiple writing styles. Polish the content for grammar, style, clarity, and flow while preserving the author's voice and intent. Improve sentence structure, word choice, and overall readability. Ensure consistency in tone and style throughout.",
		},
	}
}

// listSharedContent lists all available shared content files
func listSharedContent(sharedDir string) {
	fmt.Println("Available shared content files:")

	if _, err := os.Stat(sharedDir); os.IsNotExist(err) {
		fmt.Printf("❌ Shared directory '%s' does not exist\n", sharedDir)
		return
	}

	files, err := os.ReadDir(sharedDir)
	if err != nil {
		fmt.Printf("❌ Error reading shared directory: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Printf("📁 No files found in '%s'\n", sharedDir)
		return
	}

	for _, file := range files {
		if !file.IsDir() && isTextFile(file.Name()) {
			filePath := filepath.Join(sharedDir, file.Name())

			// Try to read first few lines for preview
			content, err := readFileContent(filePath)
			if err != nil {
				fmt.Printf("📄 %s (error reading: %v)\n", file.Name(), err)
				continue
			}

			// Show first line as preview
			lines := strings.Split(content, "\n")
			preview := "No content"
			if len(lines) > 0 && strings.TrimSpace(lines[0]) != "" {
				preview = truncateString(strings.TrimSpace(lines[0]), 80)
			}

			fmt.Printf("📄 %s\n   Preview: %s\n", file.Name(), preview)
		}
	}
}

// loadSharedContent loads shared content from specified files or patterns
func loadSharedContent(includeList, sharedDir string, verbose bool) (string, error) {
	if includeList == "" {
		return "", nil
	}

	// Split comma-separated list
	patterns := strings.Split(includeList, ",")
	var contentParts []string
	var processedFiles []string

	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		// Handle special patterns
		if pattern == "*" || pattern == "all" {
			// Load all files in shared directory
			files, err := findSharedFiles(sharedDir)
			if err != nil {
				return "", fmt.Errorf("failed to find shared files: %w", err)
			}
			for _, file := range files {
				if !contains(processedFiles, file) {
					processedFiles = append(processedFiles, file)
				}
			}
		} else {
			// Handle individual files or patterns
			files, err := resolveSharedPattern(pattern, sharedDir)
			if err != nil {
				return "", fmt.Errorf("failed to resolve pattern '%s': %w", pattern, err)
			}
			for _, file := range files {
				if !contains(processedFiles, file) {
					processedFiles = append(processedFiles, file)
				}
			}
		}
	}

	// Load content from all resolved files
	for _, filePath := range processedFiles {
		if verbose {
			fmt.Printf("📖 Loading shared content: %s\n", filePath)
		}

		// Read file content
		content, err := readFileContent(filePath)
		if err != nil {
			fmt.Printf("⚠️  Warning: failed to read shared content file '%s': %v\n", filePath, err)
			continue
		}

		// Add filename header and content
		filename := filepath.Base(filePath)
		header := fmt.Sprintf("=== %s ===", filename)
		contentParts = append(contentParts, header, content)
	}

	if len(contentParts) == 0 {
		return "", nil
	}

	return strings.Join(contentParts, "\n\n"), nil
}

// resolveSharedPattern resolves a pattern to actual file paths
func resolveSharedPattern(pattern, sharedDir string) ([]string, error) {
	var files []string

	// If it's an absolute path, use it directly
	if filepath.IsAbs(pattern) {
		if _, err := os.Stat(pattern); err == nil {
			return []string{pattern}, nil
		}
		return nil, fmt.Errorf("file not found: %s", pattern)
	}

	// Construct full path
	fullPath := filepath.Join(sharedDir, pattern)

	// Check if exact file exists
	if _, err := os.Stat(fullPath); err == nil {
		return []string{fullPath}, nil
	}

	// Try to find files matching pattern (with or without extension)
	basePattern := strings.TrimSuffix(pattern, filepath.Ext(pattern))
	
	// Find all files in shared directory
	sharedFiles, err := findSharedFiles(sharedDir)
	if err != nil {
		return nil, err
	}

	// Match against pattern
	for _, file := range sharedFiles {
		filename := filepath.Base(file)
		baseFilename := strings.TrimSuffix(filename, filepath.Ext(filename))
		
		// Check for exact match, pattern match, or substring match
		if filename == pattern || 
		   baseFilename == basePattern || 
		   strings.Contains(filename, pattern) || 
		   strings.Contains(baseFilename, pattern) {
			files = append(files, file)
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found matching pattern: %s", pattern)
	}

	return files, nil
}

// findSharedFiles finds all text files in the shared directory
func findSharedFiles(sharedDir string) ([]string, error) {
	var files []string

	if _, err := os.Stat(sharedDir); os.IsNotExist(err) {
		return files, nil
	}

	err := filepath.Walk(sharedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include text files
		if isTextFile(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}