# Input Files Directory

This directory contains input files that you want to process with Claude Copywriter.

## Usage

Process files from this directory:

```bash
# Process a single file
./claude-copywriter -input inputs/my-content.txt -role copywriter

# Batch process all files in the directory
./claude-copywriter -batch inputs/ -role copywriter

# Process with shared content
./claude-copywriter -input inputs/my-content.txt -role copywriter -include "*"
```

## Supported File Types

Claude Copywriter supports various text file formats:
- `.txt` - Plain text files
- `.md` - Markdown files  
- `.html` - HTML files
- `.json` - JSON files
- `.csv` - CSV files
- And many other text-based formats

## Example Files

Place your content files here such as:
- `blog-post.md` - Blog posts to be optimized
- `marketing-copy.txt` - Marketing content to transform
- `documentation.md` - Technical docs to improve
- `email-template.html` - Email content to enhance
