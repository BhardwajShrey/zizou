package review

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/shreybhardwaj/zizou/internal/cache"
	"github.com/shreybhardwaj/zizou/internal/diff"
)

// Client defines the interface for an AI client
type Client interface {
	SendMessage(ctx context.Context, prompt string) (string, error)
}

// Reviewer performs code reviews using an AI client
type Reviewer struct {
	client Client
	cache  cache.Cache
}

// NewReviewer creates a new reviewer
func NewReviewer(client Client, cache cache.Cache) *Reviewer {
	return &Reviewer{
		client: client,
		cache:  cache,
	}
}

// Review performs a code review on the given diff
func (r *Reviewer) Review(ctx context.Context, d *diff.Diff) (*Result, error) {
	// Generate cache key from diff content
	cacheKey := generateCacheKey(d)

	// Check cache
	if cachedResult, err := r.cache.Get(cacheKey); err == nil {
		var result Result
		if err := json.Unmarshal([]byte(cachedResult), &result); err == nil {
			return &result, nil
		}
	}

	// Build prompt
	prompt := BuildPrompt(d)

	// Send to AI client
	response, err := r.client.SendMessage(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get review from AI: %w", err)
	}

	// Parse response
	result, err := parseResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Cache result
	resultJSON, err := json.Marshal(result)
	if err == nil {
		_ = r.cache.Set(cacheKey, string(resultJSON))
	}

	return result, nil
}

// generateCacheKey creates a cache key from the diff
func generateCacheKey(d *diff.Diff) string {
	// Create a stable string representation of the diff
	data, _ := json.Marshal(d)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// parseResponse parses the AI response into a Result
func parseResponse(response string) (*Result, error) {
	var result Result

	// Try to find JSON in the response (in case AI added extra text)
	startIdx := -1
	endIdx := -1

	for i, ch := range response {
		if ch == '{' && startIdx == -1 {
			startIdx = i
		}
		if ch == '}' {
			endIdx = i + 1
		}
	}

	if startIdx == -1 || endIdx == -1 {
		return nil, fmt.Errorf("no JSON found in response")
	}

	jsonStr := response[startIdx:endIdx]

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &result, nil
}
