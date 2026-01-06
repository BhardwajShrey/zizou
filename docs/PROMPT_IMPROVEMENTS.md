# System Prompt Improvements Summary

## What Changed

The system prompt for Claude has been significantly enhanced to act as a more effective senior code reviewer.

## Key Improvements

### 1. **Senior Reviewer Persona** ✨
**Before**: "You are an expert code reviewer"
**After**: "You are a senior software engineer with 15+ years of experience conducting code reviews across multiple languages and domains. Your reviews are known for being thorough yet practical..."

### 2. **Review Philosophy** 🎯
Added explicit philosophy section:
- Trust but verify
- Impact-driven prioritization
- Actionable feedback requirement
- Teaching mindset

### 3. **Detailed Severity Guidelines** 📊
**Before**: Just category names
**After**: Clear criteria for each level:
- **critical**: Security vulnerabilities, data corruption, crashes
- **high**: Likely bugs, significant performance issues
- **medium**: Potential bugs, moderate issues
- **low**: Minor issues, optimizations
- **info**: Suggestions, educational notes

### 4. **Enhanced Focus Areas** 🔍
Expanded from general categories to specific examples:

**Security**:
- Injection vulnerabilities (SQL, command, XSS, path traversal)
- Authentication/authorization bypasses
- Sensitive data exposure
- Race conditions

**Bugs**:
- Logic errors and incorrect algorithms
- Off-by-one errors and boundary conditions
- Null/nil pointer dereferences
- Resource leaks
- Error handling gaps

**Performance**:
- Unnecessary allocations
- O(n²) algorithms
- Blocking operations on hot paths
- Memory leaks

**Maintainability**:
- Overly complex logic
- Poor naming
- Missing validation
- Code duplication

### 5. **Language-Specific Guidance** 💻
Added specific checks for each language:

**Go**:
- Missing error handling
- Goroutine leaks
- Race conditions
- Improper defer usage in loops

**C++**:
- Memory leaks, double-frees
- RAII violations
- Iterator invalidation
- Undefined behavior

**Python**:
- Mutable default arguments
- Missing resource cleanup
- Exception handling issues

**JavaScript**:
- Unhandled promise rejections
- Missing null checks
- Memory leaks in event listeners

### 6. **Structured Message Format** 📝
**Before**: "Clear description of the issue"
**After**: Enforced three-part structure:

```
Issue: [Brief description of what's wrong]
Impact: [Why this matters - potential consequences]
Fix: [Specific suggestion with code example]
```

**Example provided in prompt**:
```
Issue: Potential SQL injection vulnerability via unsanitized user input.
Impact: Attackers could execute arbitrary SQL queries, potentially reading or modifying sensitive data.
Fix: Use parameterized queries instead:
  db.Query("SELECT * FROM users WHERE id = ?", userId)
Instead of:
  db.Query("SELECT * FROM users WHERE id = " + userId)
```

### 7. **Clear Rules** 📋
Expanded from 6 to 7 explicit rules:
1. Return ONLY valid JSON
2. Be specific about line numbers
3. **Suggest concrete fixes** (emphasized)
4. Focus on significant issues
5. Consider context
6. Empty is okay
7. Explain reasoning (WHY, not just WHAT)

### 8. **What NOT to Flag** 🚫
New section to reduce noise:
- Minor style inconsistencies
- Personal preferences
- Issues handled by linters
- Complex suggestions without clear benefit

## Impact

The improved prompt should:
- ✅ Produce more actionable feedback with concrete fix suggestions
- ✅ Better prioritize issues by severity and impact
- ✅ Catch language-specific bugs and patterns
- ✅ Reduce false positives and style nitpicks
- ✅ Provide educational context to help developers learn
- ✅ Generate consistent, structured responses

## Files Modified

- `internal/client/reviewer_client.go` - Updated `buildReviewPrompt()` method
- `docs/SYSTEM_PROMPT.md` - Complete prompt documentation (new)
- `docs/PROMPT_IMPROVEMENTS.md` - This summary (new)

## Testing

Build verified: ✅
```bash
go build -v
# Success - no errors
```

## Example Usage

The improved prompt is automatically used when running:
```bash
git diff | ./zizou
# or
./zizou --file changes.diff
```

Claude will now respond with more detailed, actionable feedback including specific fix suggestions.
