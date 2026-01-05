package client

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// EnhancedClient wraps the Anthropic SDK with rate limiting and retry logic
type EnhancedClient struct {
	client      *anthropic.Client
	config      *Config
	rateLimiter *RateLimiter
}

// NewEnhancedClient creates a new enhanced Claude API client
func NewEnhancedClient(config *Config) (*EnhancedClient, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create Anthropic SDK client
	client := anthropic.NewClient(
		option.WithAPIKey(config.APIKey),
		option.WithRequestTimeout(config.Timeout),
	)

	// Create rate limiter
	var rateLimiter *RateLimiter
	if config.EnableRateLimiting {
		rateLimiter = NewRateLimiter(config.RateLimit, time.Minute)
	}

	return &EnhancedClient{
		client:      &client,
		config:      config,
		rateLimiter: rateLimiter,
	}, nil
}

// SendMessage sends a message to Claude with retry logic and rate limiting
func (ec *EnhancedClient) SendMessage(ctx context.Context, prompt string) (string, error) {
	var lastErr error

	for attempt := 0; attempt <= ec.config.MaxRetries; attempt++ {
		// Wait for rate limiter
		if ec.rateLimiter != nil {
			if err := ec.rateLimiter.Wait(ctx); err != nil {
				return "", fmt.Errorf("rate limiter error: %w", err)
			}
		}

		// Try to send message
		response, err := ec.sendMessageAttempt(ctx, prompt)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// Check if we should retry
		if !ec.shouldRetry(err, attempt) {
			break
		}

		// Calculate backoff delay
		delay := ec.calculateBackoff(attempt)

		// Wait before retry
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	return "", fmt.Errorf("failed after %d retries: %w", ec.config.MaxRetries, lastErr)
}

// sendMessageAttempt performs a single attempt to send a message
func (ec *EnhancedClient) sendMessageAttempt(ctx context.Context, prompt string) (string, error) {
	message, err := ec.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(ec.config.Model),
		MaxTokens: int64(ec.config.MaxTokens),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		return "", ec.wrapError(err)
	}

	// Extract text from response
	if len(message.Content) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	// Get the text from the first content block
	textBlock := message.Content[0].Text
	return textBlock, nil
}

// shouldRetry determines if we should retry based on the error type
func (ec *EnhancedClient) shouldRetry(err error, attempt int) bool {
	if attempt >= ec.config.MaxRetries {
		return false
	}

	// Check for retryable errors
	if anthropicErr, ok := err.(*anthropic.Error); ok {
		switch anthropicErr.StatusCode {
		case http.StatusTooManyRequests: // 429
			return true
		case http.StatusServiceUnavailable: // 503
			return true
		case http.StatusGatewayTimeout: // 504
			return true
		case http.StatusInternalServerError: // 500
			return true
		default:
			// Don't retry client errors (4xx except 429)
			if anthropicErr.StatusCode >= 400 && anthropicErr.StatusCode < 500 {
				return false
			}
			return true
		}
	}

	// Retry on network errors
	return true
}

// calculateBackoff calculates exponential backoff delay
func (ec *EnhancedClient) calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: delay * 2^attempt
	delay := ec.config.RetryDelay * time.Duration(math.Pow(2, float64(attempt)))

	// Cap at max delay
	if delay > ec.config.MaxRetryDelay {
		delay = ec.config.MaxRetryDelay
	}

	return delay
}

// wrapError wraps API errors with more context
func (ec *EnhancedClient) wrapError(err error) error {
	if anthropicErr, ok := err.(*anthropic.Error); ok {
		switch anthropicErr.StatusCode {
		case http.StatusUnauthorized:
			return fmt.Errorf("authentication failed: check your API key")
		case http.StatusForbidden:
			return fmt.Errorf("access forbidden: insufficient permissions")
		case http.StatusTooManyRequests:
			return fmt.Errorf("rate limit exceeded: %w", err)
		case http.StatusInternalServerError:
			return fmt.Errorf("API server error: %w", err)
		case http.StatusServiceUnavailable:
			return fmt.Errorf("API service unavailable: %w", err)
		default:
			return fmt.Errorf("API error (status %d): %w", anthropicErr.StatusCode, err)
		}
	}

	return err
}

// GetUsage returns usage statistics (if available from last call)
func (ec *EnhancedClient) GetUsage() *UsageStats {
	// This would need to be tracked across calls
	// For now, return nil - can be enhanced later
	return nil
}

// UsageStats tracks API usage
type UsageStats struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
	RequestCount int
}
