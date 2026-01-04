# Zizou

An AI-powered code review CLI tool that uses Claude to analyze git diffs and provide structured feedback.

## Features

- Reads git diffs from stdin or a file
- Sends the diff to Claude API for intelligent analysis
- Returns structured review comments with severity levels
- Caches duplicate results to avoid unnecessary API calls
- Multiple output formats (text, JSON, markdown)

## Installation

```bash
go install github.com/shreybhardwaj/zizou@latest
```

Or build from source:

```bash
git clone https://github.com/shreybhardwaj/zizou
cd zizou
go build
```

## Usage

### Basic Usage

```bash
# Review changes from stdin
git diff | zizou

# Review changes from a file
zizou --file changes.diff

# Specify output format
git diff | zizou --output json
git diff | zizou --output markdown
```

### API Key

Set your Anthropic API key:

```bash
export ANTHROPIC_API_KEY="your-api-key"
```

Or pass it directly:

```bash
git diff | zizou --api-key "your-api-key"
```

### Cache Options

```bash
# Disable caching
git diff | zizou --no-cache

# Custom cache directory
git diff | zizou --cache-dir /path/to/cache
```

## Output Formats

### Text (default)
Human-readable text output with issue details and severity levels.

### JSON
Structured JSON output for integration with other tools:

```json
{
  "comments": [
    {
      "file": "main.go",
      "line": 42,
      "severity": "high",
      "category": "security",
      "message": "Potential SQL injection vulnerability"
    }
  ],
  "summary": "Found 1 issue requiring attention"
}
```

### Markdown
Formatted markdown output with severity grouping and icons.

## Severity Levels

- **critical**: Security vulnerabilities, data loss risks
- **high**: Bugs, performance issues, incorrect logic
- **medium**: Code quality, maintainability concerns
- **low**: Style issues, minor improvements
- **info**: Suggestions, general observations

## Categories

- `security`: Security vulnerabilities
- `performance`: Performance issues
- `bug_risk`: Potential bugs or logic errors
- `maintainability`: Code maintainability concerns
- `style`: Code style issues
- `best_practice`: Best practice violations
- `other`: Other issues

## Project Structure

```
zizou/
├── main.go                    # Entry point
├── cmd/
│   └── root.go                # CLI command setup
├── internal/
│   ├── diff/                  # Git diff parsing
│   ├── client/                # Claude API client
│   ├── review/                # Review orchestration
│   ├── cache/                 # Result caching
│   └── output/                # Output formatting
```

## Examples

```bash
# Review uncommitted changes
git diff | zizou

# Review staged changes
git diff --staged | zizou --output markdown

# Review specific commit
git show HEAD | zizou

# Review changes between branches
git diff main..feature-branch | zizou --output json
```

## License

MIT
