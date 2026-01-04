package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

// Cache defines the interface for caching review results
type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

// FileCache implements file-based caching
type FileCache struct {
	dir string
}

// NewFileCache creates a new file-based cache
func NewFileCache(dir string) (*FileCache, error) {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &FileCache{
		dir: dir,
	}, nil
}

// Get retrieves a value from the cache
func (c *FileCache) Get(key string) (string, error) {
	path := c.getPath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Set stores a value in the cache
func (c *FileCache) Set(key string, value string) error {
	path := c.getPath(key)
	return os.WriteFile(path, []byte(value), 0644)
}

// getPath returns the file path for a cache key
func (c *FileCache) getPath(key string) string {
	return filepath.Join(c.dir, key+".json")
}

// NoOpCache is a cache implementation that doesn't cache anything
type NoOpCache struct{}

// NewNoOpCache creates a new no-op cache
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// Get always returns an error (cache miss)
func (c *NoOpCache) Get(key string) (string, error) {
	return "", fmt.Errorf("cache disabled")
}

// Set does nothing
func (c *NoOpCache) Set(key string, value string) error {
	return nil
}
