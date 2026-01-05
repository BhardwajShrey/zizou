package client

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Model != "claude-sonnet-4-20250514" {
		t.Errorf("Model = %s, want claude-sonnet-4-20250514", config.Model)
	}

	if config.MaxTokens != 4096 {
		t.Errorf("MaxTokens = %d, want 4096", config.MaxTokens)
	}

	if config.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want 60s", config.Timeout)
	}

	if config.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", config.MaxRetries)
	}

	if !config.EnableRateLimiting {
		t.Error("EnableRateLimiting = false, want true")
	}
}

func TestConfig_LoadFromEnv(t *testing.T) {
	// Set test environment variable
	os.Setenv("ANTHROPIC_API_KEY_ZIZOU", "test-api-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY_ZIZOU")

	config := DefaultConfig()
	err := config.LoadFromEnv()

	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v, want nil", err)
	}

	if config.APIKey != "test-api-key" {
		t.Errorf("APIKey = %s, want test-api-key", config.APIKey)
	}
}

func TestConfig_LoadFromEnv_AlternativeKey(t *testing.T) {
	// Test CLAUDE_API_KEY as alternative
	os.Setenv("CLAUDE_API_KEY", "alternative-key")
	defer os.Unsetenv("CLAUDE_API_KEY")

	config := DefaultConfig()
	err := config.LoadFromEnv()

	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v, want nil", err)
	}

	if config.APIKey != "alternative-key" {
		t.Errorf("APIKey = %s, want alternative-key", config.APIKey)
	}
}

func TestConfig_LoadFromEnv_NoKey(t *testing.T) {
	// Ensure no env vars are set
	os.Unsetenv("ANTHROPIC_API_KEY_ZIZOU")
	os.Unsetenv("CLAUDE_API_KEY")

	config := DefaultConfig()
	err := config.LoadFromEnv()

	if err == nil {
		t.Error("LoadFromEnv() error = nil, want error when API key not set")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name: "valid config",
			config: &Config{
				APIKey:             "test-key",
				Model:              "claude-sonnet-4",
				MaxTokens:          1000,
				Timeout:            30 * time.Second,
				MaxRetries:         3,
				EnableRateLimiting: true,
				RateLimit:          50,
			},
			wantError: false,
		},
		{
			name: "missing API key",
			config: &Config{
				Model:      "claude-sonnet-4",
				MaxTokens:  1000,
				Timeout:    30 * time.Second,
				MaxRetries: 3,
				RateLimit:  50,
			},
			wantError: true,
		},
		{
			name: "missing model",
			config: &Config{
				APIKey:     "test-key",
				MaxTokens:  1000,
				Timeout:    30 * time.Second,
				MaxRetries: 3,
				RateLimit:  50,
			},
			wantError: true,
		},
		{
			name: "invalid max tokens",
			config: &Config{
				APIKey:     "test-key",
				Model:      "claude-sonnet-4",
				MaxTokens:  0,
				Timeout:    30 * time.Second,
				MaxRetries: 3,
				RateLimit:  50,
			},
			wantError: true,
		},
		{
			name: "invalid timeout",
			config: &Config{
				APIKey:     "test-key",
				Model:      "claude-sonnet-4",
				MaxTokens:  1000,
				Timeout:    0,
				MaxRetries: 3,
				RateLimit:  50,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestNewConfigFromEnv(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY_ZIZOU", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY_ZIZOU")

	config, err := NewConfigFromEnv()

	if err != nil {
		t.Fatalf("NewConfigFromEnv() error = %v, want nil", err)
	}

	if config.APIKey != "test-key" {
		t.Errorf("APIKey = %s, want test-key", config.APIKey)
	}

	// Should have default values
	if config.Model == "" {
		t.Error("Model should have default value")
	}
}
