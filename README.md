<div align="center">

<img src="https://github.com/user-attachments/assets/18e50d59-3c1f-4577-8c27-372dcecc244e" width="300" height="300">

# Claude Copywriting System
</div>

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

```text
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

Set your Anthropic API key either:

```bash
export ANTHROPIC_API_KEY='your-api-key-here'
```

Or add it to `config.json`:

```json
{
  "api_key": "your-api-key-here",
  "base_url": "https://api.anthropic.com",
  "model": "claude-3-sonnet-20240229"
}
```

### Basic Usage

```bash
./claude-copywriter -input myfile.txt -role copywriter
```

### Batch Processing

```bash
./claude-copywriter -batch inputs/ -role copywriter -concurrent 5
```

### Shared Content Integration

```bash
# Include all shared content
./claude-copywriter -input myfile.txt -role copywriter -include "*"

# Include specific files
./claude-copywriter -input myfile.txt -role copywriter -include "brand-voice,company-info"
```

### Custom Prompts

```bash
./claude-copywriter -input myfile.txt -role copywriter -prompt "Make this more technical and include code examples"
```

## Advanced Usage

### Managing Roles

```bash
# List all available roles
./claude-copywriter -list-roles

# View shared content files
./claude-copywriter -list-shared
```

### Setting Up Shared Content

```bash
# Create shared content directory
mkdir shared

# Add brand guidelines
echo "Your brand voice guidelines here" > shared/brand-voice.md
```

### Using Shared Content

```bash
# Include all shared files
./claude-copywriter -input content.txt -role copywriter -include "*"

# Include files matching pattern
./claude-copywriter -input content.txt -role copywriter -include "brand"

# Include specific files
./claude-copywriter -input content.txt -role copywriter -include "brand-voice,style-guide"
```

### Content Integration

When you include shared content, it gets combined with your main content like this:

```text
=== brand-voice.md ===
[Your brand voice content]

=== company-info.md ===
[Your company information]

---

MAIN CONTENT TO PROCESS:
[Your original file content]
```

## Available Roles

- **copywriter**: Expert copywriting with AIDA, PAS frameworks
- **researcher**: Detailed research insights and fact-checking
- **reviewer**: Critical content review and improvement suggestions
- **formatter**: Structure and readability optimization
- **summarizer**: Concise summaries of complex content
- **editor**: Grammar, style, and flow improvements
- **seo-optimizer**: Search engine optimization
- **social-media**: Platform-specific social content
- **technical-writer**: Clear technical documentation
- **blog-writer**: Engaging SEO-friendly blog posts
- **email-marketer**: Effective email campaigns
- **academic-writer**: Scholarly writing standards
- **creative-writer**: Storytelling and narrative techniques
- **analyst**: Business analysis and strategic insights
- **translator**: Professional translation and localization
- **natural**: Human-like writing without AI-typical phrases

## Tips & Best Practices

- **Rate Limiting**: Built-in 100ms delay between concurrent requests
- **File Types**: Supports .txt, .md, .html, .json, .csv, and more
- **Output**: Files are saved as `input_role.extension`
- **Verbose Mode**: Use `-verbose` for detailed processing info
- **Token Tracking**: Monitor API usage with automatic reporting

## Troubleshooting

### Common Issues

1. **API Key Error**: Ensure your API key is set correctly
2. **File Not Found**: Check file paths and permissions
3. **Rate Limits**: Reduce concurrent requests with `-concurrent 1`

### Support

For issues with the Claude API, visit: <https://docs.anthropic.com>

## License

MIT License - feel free to modify and distribute.
