# Running Zizou - Complete Guide

## Your Diff File

We've saved your diff to `test_diff.txt`. Here's how to use Zizou with it:

## Prerequisites

Set your Claude API key:
```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

Or pass it via flag:
```bash
--api-key "your-api-key-here"
```

---

## Method 1: Read Diff from File

```bash
# Text output (default)
./zizou --file test_diff.txt

# JSON output
./zizou --file test_diff.txt --output json

# Markdown output
./zizou --file test_diff.txt --output markdown

# With explicit API key
./zizou --file test_diff.txt --api-key "sk-ant-..."
```

---

## Method 2: Pipe from Git (Most Common)

```bash
# Review uncommitted changes
git diff | ./zizou

# Review staged changes
git diff --staged | ./zizou

# Review specific commit
git show HEAD | ./zizou

# Review between branches
git diff main..feature-branch | ./zizou

# Review last commit
git diff HEAD~1 HEAD | ./zizou
```

---

## Method 3: From Stdin

```bash
# Pipe your diff directly
cat test_diff.txt | ./zizou

# Or use input redirection
./zizou < test_diff.txt
```

---

## Output Formats

### Text (Default)
```bash
./zizou --file test_diff.txt --output text
```

**Example output:**
```
Code Review Results
============================================================

[1] app/main.go:25
    Severity: high | Category: bug_risk
    The deferred os.Remove(jailPath) could fail silently. Consider
    checking the error or logging it.

[2] app/main.go:40
    Severity: info | Category: best_practice
    Good refactoring! Returning the path allows proper cleanup.

------------------------------------------------------------
Summary:
Found 2 issues requiring attention. The change improves resource
management by returning the temp directory path for cleanup.
```

### JSON
```bash
./zizou --file test_diff.txt --output json
```

**Example output:**
```json
{
  "comments": [
    {
      "file": "app/main.go",
      "line": 25,
      "severity": "high",
      "category": "bug_risk",
      "message": "Deferred os.Remove could fail silently...",
      "snippet": "defer os.Remove(jailPath)"
    }
  ],
  "summary": "Found 2 issues requiring attention..."
}
```

### Markdown
```bash
./zizou --file test_diff.txt --output markdown
```

**Example output:**
```markdown
# Code Review Results

## 🔴 Critical (0)

## 🟠 High (1)

### app/main.go:25
**Category:** bug_risk

The deferred os.Remove(jailPath) could fail silently...

## Summary

Found 2 issues requiring attention...
```

---

## Cache Options

```bash
# Disable caching
./zizou --file test_diff.txt --no-cache

# Custom cache directory
./zizou --file test_diff.txt --cache-dir /tmp/zizou-cache

# Default cache location
# ~/.zizou/cache
```

---

## Complete Examples

### Example 1: Review Your Diff
```bash
./zizou --file test_diff.txt --output markdown
```

### Example 2: CI/CD Integration
```bash
git diff origin/main...HEAD | ./zizou --output json > review.json
```

### Example 3: Pre-commit Hook
```bash
git diff --staged | ./zizou --output text
```

### Example 4: Review PR
```bash
gh pr diff 123 | ./zizou --output markdown > pr_review.md
```

---

## What Zizou Will Analyze

From your diff, Zizou will review:

✅ **Hunk 1** (Lines 23-30):
- Changed function call to capture return value
- Added defer cleanup

✅ **Hunk 2** (Lines 37-50):
- Changed function signature to return string
- Removed premature cleanup

✅ **Hunk 3** (Lines 58-63):
- Added return statement

### Potential Issues Zizou Might Flag:

1. **Error handling** - `defer os.Remove(jailPath)` doesn't check errors
2. **Resource leaks** - If the function panics before return
3. **Security** - `os.Chmod(tempDirPath, 0777)` is very permissive
4. **Best practices** - Overall good refactoring pattern

---

## Quick Test (Without API Call)

Since Zizou needs a Claude API key, here's how to test the parser without API:

```bash
# Just parse and show statistics
go run test_parser.go
```

This will show what Zizou extracts from your diff without calling the API.

---

## Getting Help

```bash
./zizou --help
```

**Output:**
```
Zizou is a Go-based code review CLI tool that:
- Reads git diffs from stdin or a file
- Sends the diff to Claude API for analysis
- Returns structured review comments with severity levels
- Caches duplicate results to avoid API calls

Usage:
  zizou [flags]

Flags:
      --api-key string     Claude API key (or set ANTHROPIC_API_KEY env var)
      --cache-dir string   Cache directory (defaults to ~/.zizou/cache)
  -f, --file string        Input file containing git diff (defaults to stdin)
  -h, --help               help for zizou
      --no-cache           Disable caching
  -o, --output string      Output format (text, json, markdown) (default "text")
```

---

## Pro Tips

1. **Use with git aliases:**
   ```bash
   git config --global alias.review '!git diff | zizou'
   git review  # Review your changes
   ```

2. **Combine with other tools:**
   ```bash
   git diff | zizou --output json | jq '.comments[] | select(.severity == "critical")'
   ```

3. **CI/CD Integration:**
   ```yaml
   - name: Review PR
     run: |
       git diff origin/main...HEAD | \
         ./zizou --output json > review.json
   ```

4. **Save reviews:**
   ```bash
   git diff | ./zizou --output markdown > review_$(date +%Y%m%d).md
   ```

---

## Troubleshooting

### "API key required"
Set your API key:
```bash
export ANTHROPIC_API_KEY="sk-ant-..."
```

### "Failed to parse diff"
Check your diff is valid:
```bash
cat test_diff.txt  # Should show valid git diff format
```

### "Empty diff content"
Make sure the diff file is not empty:
```bash
wc -l test_diff.txt  # Should show line count
```

---

## Next Steps

1. Get a Claude API key from https://console.anthropic.com
2. Set the environment variable
3. Run Zizou on your diff
4. Review the output and fix issues
5. Integrate into your workflow
