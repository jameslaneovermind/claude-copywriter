# Claude Copywriting System

A powerful Go-based CLI tool that processes text files through Claude AI with different roles and research capabilities.

## Features

- 🎭 **Multiple Roles**: Copywriter, Researcher, Reviewer, Formatter, Summarizer, Editor, and more
- 🔍 **AI Research**: Deep research integration via Perplexity API with real-time data and citations
- 📁 **Batch Processing**: Process entire directories of files
- 🚀 **Concurrent Processing**: Configurable concurrent API calls
- 📖 **Shared Content**: Include brand guidelines, company info, and other context
- 🔒 **Privacy-Aware**: Sensitive data directories excluded from version control

## Quick Start

1. **Install and build**:

   ```bash
   git clone <your-repo>
   cd claude-copywriter
   go build -o claude-copywriter
   ```

2. **Set API keys**:

   ```bash
   export ANTHROPIC_API_KEY='your-api-key-here'
   export PERPLEXITY_API_KEY='your-perplexity-key'  # Optional, for research
   ```

3. **Basic usage**:

   ```bash
   ./claude-copywriter -input document.md -role copywriter
   ```

## Directory Structure

```
claude-copywriter/
├── main.go, client.go, processor.go    # Core application
├── config.json                         # Configuration and roles
├── shared/                              # Brand content (gitignored)
├── inputs/                              # Input files (gitignored)
├── shared-example/                      # Example structures
└── tests/                               # Test files
```

## Core Usage

### Basic Processing

```bash
# Transform content with different roles
./claude-copywriter -input article.md -role copywriter
./claude-copywriter -input docs.md -role technical-writer
./claude-copywriter -input content.md -role editor
```

### Research Integration

```bash
# Research-only mode (saves research to file)
./claude-copywriter -input article.md -research "DevOps trends 2024" -research-only

# Research + processing (includes research in context)
./claude-copywriter -input article.md -role copywriter -research "cloud security trends"
```

### Shared Content

```bash
# Include brand guidelines and context
./claude-copywriter -input content.md -role copywriter -include "brand-voice,company-info"

# Include all shared content
./claude-copywriter -input content.md -role copywriter -include "*"

# Pattern matching
./claude-copywriter -input content.md -role copywriter -include "brand"
```

### Batch Processing

```bash
# Process all files in a directory
./claude-copywriter -batch inputs/ -role copywriter -concurrent 3
```

## Research Capabilities

The Perplexity API integration provides:

- **Real-time information**: Current web content and trends
- **Authoritative sources**: Citations from reliable publications
- **Technical depth**: Latest API docs and best practices

### Example Research Workflow

```bash
# 1. Research current trends
./claude-copywriter -input comparison.md -research "infrastructure scanning tools 2024" -research-only

# 2. Create enhanced content with research
./claude-copywriter -input comparison.md -role copywriter -research "TFLint vs security tools" -include "brand-voice"

# 3. Technical deep dive
./claude-copywriter -input comparison.md -role technical-writer -research "Terraform security best practices"
```

## Configuration

### Available Roles

```bash
./claude-copywriter -list-roles
```

Built-in roles include: `copywriter`, `researcher`, `reviewer`, `formatter`, `summarizer`, `editor`, `seo-optimizer`, `social-media`, `technical-writer`, `blog-writer`, `email-marketer`, `academic-writer`, `creative-writer`, `analyst`, `translator`, `natural`.

### Custom Configuration

Edit `config.json` to:

- Add custom roles and prompts
- Configure API endpoints and models
- Set default parameters

### Shared Content Setup

```bash
# List available shared content
./claude-copywriter -list-shared

# Create shared content files
mkdir -p shared
echo "# Brand Voice Guidelines..." > shared/brand-voice.md
echo "# Company Information..." > shared/company-info.md
```

## Security & Privacy

🔒 **Important**: `shared/` and `inputs/` directories are excluded from version control by default.

**Best Practices**:

- Keep sensitive content in excluded directories
- Use environment variables for API keys
- Never commit real API keys to version control

## API Key Setup

### Claude API (Required)

1. Get API key from [Anthropic Console](https://console.anthropic.com/)
2. Set: `export ANTHROPIC_API_KEY='your-key'`

### Perplexity API (Optional - for research)

1. Get API key from [Perplexity AI](https://www.perplexity.ai/)
2. Set: `export PERPLEXITY_API_KEY='your-key'`

## Advanced Usage

### Custom Prompts

```bash
./claude-copywriter -input content.md -role copywriter -prompt "Focus on technical accuracy and include code examples"
```

### Output Control

```bash
./claude-copywriter -input article.md -role copywriter -output custom-name.md
```

### Verbose Mode

```bash
./claude-copywriter -input content.md -role copywriter -verbose
```

## Common Use Cases

- **Content Enhancement**: Add current trends and data to existing articles
- **Technical Documentation**: Create comprehensive guides with latest best practices  
- **Marketing Copy**: Transform content with brand voice consistency
- **Research Reports**: Gather and synthesize information from multiple sources
- **Competitive Analysis**: Research industry trends and positioning
- **Blog Posts**: Create SEO-optimized content with proper structure

## Troubleshooting

- **API Errors**: Ensure API keys are set correctly
- **File Not Found**: Check input file paths and permissions
- **Research Failures**: Verify Perplexity API key and network connectivity
- **Build Issues**: Ensure Go 1.21+ is installed

For more help: `./claude-copywriter -h`
