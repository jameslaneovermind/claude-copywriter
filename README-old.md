// go.mod
module claude-copywriter

go 1.21

// This project uses only the Go standard library
// No external dependencies required!

//---

// README.md
# Claude Copywriting System

A powerful Go-based CLI tool that processes text files through Claude AI with different roles and prompts.

## Features

- 🎭 **Multiple Roles**: Copywriter, Researcher, Reviewer, Formatter, Summarizer, Editor, and more
- 📁 **Batch Processing**: Process entire directories of files
- 🚀 **Concurrent Processing**: Configurable concurrent API calls
- 🔧 **Configurable**: Easy-to-modify JSON configuration
- 📊 **Token Tracking**: Monitor API usage
- 🎯 **Flexible Output**: Custom output filenames and locations
- 📖 **Shared Content**: Include brand guidelines, company info, and other context
- 🔒 **Privacy-Aware**: Sensitive data directories excluded from version control

## Quick Start

1. **Set up the environment**:
```bash
git clone <your-repo>
cd claude-copywriter
chmod +x setup.sh
./setup.sh
```

2. **Set your API key**:
```bash
export ANTHROPIC_API_KEY='your-api-key-here'
```

3. **Try it out**:
```bash
./claude-copywriter -list-roles
echo "Test content" > test.txt
./claude-copywriter -input test.txt -role copywriter
```

## Directory Structure

```
claude-copywriter/
├── main.go              # Main application
├── client.go            # Claude API client  
├── processor.go         # File processing logic
├── config.json          # Configuration and roles
├── shared/              # Shared content (gitignored)
│   ├── brand-voice.md   # Brand guidelines
│   ├── company-info.md  # Company context
│   └── ...              # Other shared content
├── inputs/              # Input files (gitignored)  
│   ├── blog-post.md     # Content to process
│   └── ...              # Other input files
├── shared-example/      # Example shared content structure
├── inputs-example/      # Example input files structure
└── outputs/             # Generated output files
```

## Security & Privacy

🔒 **Important**: The `shared/` and `inputs/` directories are excluded from version control by default to protect sensitive data like:
- Brand guidelines and messaging
- Company confidential information  
- Customer content and data
- API keys and credentials

**Best Practices**:
- Keep sensitive content in `shared/` and `inputs/` directories
- Use `shared-example/` and `inputs-example/` for documentation
- Never commit real API keys to version control
- Use environment variables for API keys

## Installation

1. **Clone or create the project**:
```bash
mkdir claude-copywriter
cd claude-copywriter
```

2. **Initialize Go module**:
```bash
go mod init claude-copywriter
```

3. **Create the main files** (copy the provided Go code into these files):
- `main.go` - Main application logic
- `client.go` - Claude API client
- `processor.go` - File processing logic

4. **Build the application**:
```bash
go build -o claude-copywriter
```

## Configuration

### API Key Setup

Set your Anthropic API key as an environment variable:
```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

Or add it to the `config.json` file (created automatically on first run).

### Configuration File

The application creates a `config.json` file with default settings:

```json
{
  "api_key": "",
  "base_url": "https://api.anthropic.com",
  "model": "claude-sonnet-4-20250514",
  "roles": {
    "copywriter": "You are an expert copywriter...",
    "researcher": "You are a thorough researcher...",
    "reviewer": "You are a critical reviewer...",
    "formatter": "You are a formatting specialist...",
    "summarizer": "You are a summarization expert...",
    "editor": "You are a professional editor..."
  }
}
```

## Usage Examples

### Basic Usage
```bash
# Process a single file with default copywriter role
./claude-copywriter -input article.txt

# Use a specific role
./claude-copywriter -input article.txt -role researcher

# Add custom instructions for this specific request
./claude-copywriter -input article.txt -role copywriter -prompt "Focus on ROI and include specific metrics"

# Specify output file
./claude-copywriter -input article.txt -output research_report.txt -role researcher
```

### Batch Processing
```bash
# Process all text files in a directory
./claude-copywriter -batch ./documents -role editor

# With custom concurrency
./claude-copywriter -batch ./documents -role copywriter -concurrent 5
```

### Shared Content Integration
```bash
# Include company messaging (creates shared/overmind.md automatically)
./claude-copywriter -input article.txt -role copywriter -include overmind

# Include multiple shared files
./claude-copywriter -input proposal.txt -role copywriter -include overmind,company-info

# Combine shared content with custom instructions
./claude-copywriter -input blog-draft.txt -role copywriter -include overmind -prompt "Write for technical audience, include code examples"

# List available shared content
./claude-copywriter -list-shared

# Batch processing with shared content
./claude-copywriter -batch ./drafts -role blog-writer -include overmind
```

### Custom Prompts
```bash
# Add specific instructions for one-off requests
./claude-copywriter -input features.txt -role copywriter -prompt "Focus on enterprise customers, emphasize security and compliance"

# Combine role, shared content, and custom instructions
./claude-copywriter -input article.txt -role blog-writer -include overmind,brand-voice -prompt "Target CTO audience, include technical depth"

# Different instructions for batch processing
./claude-copywriter -batch ./press-releases -role copywriter -prompt "Keep under 300 words, include strong headlines"
```

### Managing Roles
```bash
# List all available roles
./claude-copywriter -list-roles

# Verbose output for debugging
./claude-copywriter -input file.txt -role reviewer -verbose
```

## Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-input` | Input file path | Required for single file |
| `-output` | Output file path | `input_[role].ext` |
| `-role` | Role/prompt to use | `copywriter` |
| `-prompt` | Additional custom instructions | - |
| `-config` | Configuration file path | `config.json` |
| `-list-roles` | List available roles | `false` |
| `-batch` | Process directory of files | - |
| `-concurrent` | Concurrent API calls | `3` |
| `-verbose` | Verbose output | `false` |
| `-include` | Include shared content files | - |
| `-shared-dir` | Shared content directory | `shared` |
| `-list-shared` | List shared content files | `false` |

## Supported File Types

- Text files (`.txt`, `.md`, `.text`)
- Code files (`.go`, `.py`, `.js`, `.java`, etc.)
- Data files (`.csv`, `.json`, `.xml`, `.yaml`)
- Web files (`.html`, `.css`)
- Configuration files (`.ini`, `.cfg`, `.conf`)
- Log files (`.log`)

## Shared Content Management

The system supports reusable content blocks for consistent brand messaging:

### Setting Up Shared Content
```bash
# Create shared content directory
mkdir shared

# Add your company information
echo "# Overmind Product Info
Overmind is an AI-powered platform that..." > shared/overmind.md

echo "# Company Background  
Founded in 2023, our mission is..." > shared/company-info.md
```

### Using Shared Content
```bash
# Include single shared file
./claude-copywriter -input blog.txt -role copywriter -include overmind

# Include multiple files (comma-separated)
./claude-copywriter -input proposal.txt -include overmind,company-info,brand-voice

# List what's available
./claude-copywriter -list-shared
```

### Content Integration
Shared content is automatically prepended as context:
```
SHARED CONTEXT INFORMATION:
[Your overmind.md content here]

---

MAIN CONTENT TO PROCESS:
[Your input file content here]
```

This ensures Claude has your company context when writing, while keeping your source files clean.

## Customizing Roles

Edit the `config.json` file to add or modify roles:

```json
{
  "roles": {
    "blog-writer": "You are a blog writing specialist. Transform content into engaging blog posts with SEO optimization...",
    "technical-writer": "You are a technical documentation expert. Create clear, comprehensive technical documentation...",
    "social-media": "You are a social media expert. Transform content into engaging social media posts..."
  }
}
```

## Rate Limiting

The application includes built-in rate limiting:
- 100ms delay between concurrent requests
- Configurable concurrency (default: 3 concurrent requests)
- 2-minute timeout per API request

## Error Handling

- Comprehensive error reporting
- Automatic retry logic (planned for future versions)
- Detailed token usage tracking
- Batch processing continues even if individual files fail

## Development

To extend the application:

1. **Add new file types**: Modify `isTextFile()` in `processor.go`
2. **Custom output formats**: Extend the response processing in `client.go`
3. **New features**: The modular design makes it easy to add features

## License

Open source - modify and distribute as needed.

## Support

For issues with the Claude API, visit: https://docs.anthropic.com