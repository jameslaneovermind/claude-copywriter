# Shared Content Directory

This directory contains shared content files that can be included with your content processing requests.

## Usage

Use the `-include` flag to include shared content:

```bash
# Include all shared files
./claude-copywriter -input myfile.txt -role copywriter -include "*"

# Include specific files
./claude-copywriter -input myfile.txt -role copywriter -include "brand-voice,company-info"

# Include files matching a pattern
./claude-copywriter -input myfile.txt -role copywriter -include "brand"
```

## Example Files

Create markdown files in this directory such as:
- `brand-voice.md` - Your brand voice guidelines
- `company-info.md` - Company background and context
- `messaging.md` - Key messaging and positioning
- `style-guide.md` - Writing style guidelines

Each file will be included with a header showing the filename, making it easy to reference in your prompts.
