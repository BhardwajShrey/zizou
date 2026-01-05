package client

import (
	"fmt"
	"os"
	"time"
)

// Config holds configuration for the Claude API client
type Config struct {
	// API key for authentication
	APIKey string

	// Model to use (default: claude-sonnet-4-20250514)
	Model string

	// Maximum tokens in response (default: 4096)
	MaxTokens int

	// Request timeout (default: 60s)
	Timeout time.Duration

	// Maximum retry attempts (default: 3)
	MaxRetries int

	// Initial retry delay (default: 1s)
	RetryDelay time.Duration

	// Maximum retry delay (default: 30s)
	MaxRetryDelay time.Duration

	// Enable rate limiting (default: true)
	EnableRateLimiting bool

	// Requests per minute limit (default: 60)
	RateLimit int
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		APIKey:             "",
		Model:              "claude-sonnet-4-20250514",
		MaxTokens:          4096,
		Timeout:            60 * time.Second,
		MaxRetries:         3,
		RetryDelay:         1 * time.Second,
		MaxRetryDelay:      30 * time.Second,
		EnableRateLimiting: true,
		RateLimit:          50, // Conservative: 50 requests per minute
	}
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() error {
	// API Key (required)
	if key := os.Getenv("ANTHROPIC_API_KEY_ZIZOU"); key != "" {
		c.APIKey = key
	}

	// Alternative: CLAUDE_API_KEY
	if key := os.Getenv("CLAUDE_API_KEY"); key != "" {
		c.APIKey = key
	}

	if c.APIKey == "" {
		return fmt.Errorf("API key not found: set ANTHROPIC_API_KEY_ZIZOU or CLAUDE_API_KEY environment variable")
	}

	// Model (optional)
	if model := os.Getenv("CLAUDE_MODEL"); model != "" {
		c.Model = model
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	if c.Model == "" {
		return fmt.Errorf("model is required")
	}

	if c.MaxTokens <= 0 {
		return fmt.Errorf("max tokens must be positive")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	if c.RateLimit <= 0 {
		return fmt.Errorf("rate limit must be positive")
	}

	return nil
}

// NewConfigFromEnv creates a new configuration loaded from environment variables
func NewConfigFromEnv() (*Config, error) {
	config := DefaultConfig()
	if err := config.LoadFromEnv(); err != nil {
		return nil, err
	}
	return config, config.Validate()
}
