package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type ModeDuration struct {
	Work       time.Duration `json:"work"`
	ShortBreak time.Duration `json:"short_break"`
	LongBreak  time.Duration `json:"long_break"`
	Walk       time.Duration `json:"walk"`
	Water      time.Duration `json:"water"`
	Video      time.Duration `json:"video"`
}

type Config struct {
	Durations    ModeDuration `json:"durations"`
	SoundEnabled bool         `json:"sound_enabled"`
	NotifyEnabled bool        `json:"notify_enabled"`
}

func DefaultConfig() Config {
	return Config{
		Durations: ModeDuration{
			Work:       25 * time.Minute,
			ShortBreak: 5 * time.Minute,
			LongBreak:  15 * time.Minute,
			Walk:       10 * time.Minute,
			Water:      2 * time.Minute,
			Video:      20 * time.Minute,
		},
		SoundEnabled:  true,
		NotifyEnabled: true,
	}
}

func Load() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return DefaultConfig(), err
	}

	configPath := filepath.Join(homeDir, ".workpulse", "config.json")

	file, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return DefaultConfig(), err
	}

	return config, nil
}

func Save(config Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(homeDir, ".workpulse")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(dirPath, "config.json")

	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, file, 0644)
}
