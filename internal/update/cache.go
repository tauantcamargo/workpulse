package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const cacheTTL = 24 * time.Hour

// CacheEntry stores the cached update check result
type CacheEntry struct {
	LatestVersion string    `json:"latest_version"`
	CheckedAt     time.Time `json:"checked_at"`
}

// Cache handles caching of update check results
type Cache struct {
	path string
}

// NewCache creates a new cache instance
func NewCache() (*Cache, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dirPath := filepath.Join(homeDir, ".workpulse")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, err
	}

	return &Cache{
		path: filepath.Join(dirPath, "update_cache.json"),
	}, nil
}

// Get retrieves the cached version if it's still valid
func (c *Cache) Get() (string, bool) {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return "", false
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", false
	}

	if time.Since(entry.CheckedAt) > cacheTTL {
		return "", false
	}

	return entry.LatestVersion, true
}

// Set stores a version in the cache
func (c *Cache) Set(version string) error {
	entry := CacheEntry{
		LatestVersion: version,
		CheckedAt:     time.Now(),
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.path, data, 0644)
}

// Clear removes the cache file
func (c *Cache) Clear() error {
	err := os.Remove(c.path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
