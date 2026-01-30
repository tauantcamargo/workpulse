package update

import (
	"encoding/json"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	githubAPIURL = "https://api.github.com/repos/tauantcamargo/workpulse/releases/latest"
	httpTimeout  = 5 * time.Second
)

const installScriptURL = "https://raw.githubusercontent.com/tauantcamargo/workpulse/main/install.sh"

// UpdateInfo contains information about an available update
type UpdateInfo struct {
	CurrentVersion Version
	LatestVersion  Version
	UpdateURL      string
	HasUpdate      bool
}

// UpdateCheckMsg is sent when an update check completes
type UpdateCheckMsg struct {
	Info  UpdateInfo
	Error error
}

// Checker checks for new versions
type Checker struct {
	currentVersion string
	cache          *Cache
}

// NewChecker creates a new update checker
func NewChecker(currentVersion string) *Checker {
	cache, _ := NewCache()
	return &Checker{
		currentVersion: currentVersion,
		cache:          cache,
	}
}

// Check performs an update check and returns a tea.Cmd
func (c *Checker) Check() tea.Cmd {
	return func() tea.Msg {
		info, err := c.checkForUpdate()
		return UpdateCheckMsg{Info: info, Error: err}
	}
}

// CheckSync performs a synchronous update check
func (c *Checker) CheckSync() (UpdateInfo, error) {
	return c.checkForUpdate()
}

func (c *Checker) checkForUpdate() (UpdateInfo, error) {
	current := ParseVersion(c.currentVersion)

	info := UpdateInfo{
		CurrentVersion: current,
		UpdateURL:      installScriptURL,
	}

	if c.cache != nil {
		if cachedVersion, ok := c.cache.Get(); ok {
			latest := ParseVersion(cachedVersion)
			info.LatestVersion = latest
			info.HasUpdate = c.shouldUpdate(current, latest)
			return info, nil
		}
	}

	latest, err := c.fetchLatestVersion()
	if err != nil {
		return info, err
	}

	if c.cache != nil {
		c.cache.Set(latest.String())
	}

	info.LatestVersion = latest
	info.HasUpdate = c.shouldUpdate(current, latest)

	return info, nil
}

func (c *Checker) shouldUpdate(current, latest Version) bool {
	if latest.IsZero() {
		return false
	}

	if current.IsDev() {
		return true
	}

	return current.LessThan(latest)
}

func (c *Checker) fetchLatestVersion() (Version, error) {
	client := &http.Client{Timeout: httpTimeout}

	resp, err := client.Get(githubAPIURL)
	if err != nil {
		return Version{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Version{}, nil
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return Version{}, err
	}

	return ParseVersion(release.TagName), nil
}

// ClearCache clears the update cache
func (c *Checker) ClearCache() error {
	if c.cache == nil {
		return nil
	}
	return c.cache.Clear()
}
