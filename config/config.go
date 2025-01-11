package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Storage struct {
		FilePath  string `yaml:"filePath"`
		BackupDir string `yaml:"backupDir"`
	} `yaml:"storage"`

	Task struct {
		MaxTitleLength       int           `yaml:"maxTitleLength"`
		MaxDescriptionLength int           `yaml:"maxDescriptionLength"`
		DateFormat           string        `yaml:"dateFormat"`
		AutoBackup           bool          `yaml:"autoBackup"`
		BackupInterval       time.Duration `yaml:"backupInterval"`
	} `yaml:"task"`
}

var DefaultConfig = Config{
	Storage: struct {
		FilePath  string `yaml:"filePath"`
		BackupDir string `yaml:"backupDir"`
	}{
		FilePath:  "tasks.json",
		BackupDir: "backups",
	},
	Task: struct {
		MaxTitleLength       int           `yaml:"maxTitleLength"`
		MaxDescriptionLength int           `yaml:"maxDescriptionLength"`
		DateFormat           string        `yaml:"dateFormat"`
		AutoBackup           bool          `yaml:"autoBackup"`
		BackupInterval       time.Duration `yaml:"backupInterval"`
	}{
		MaxTitleLength:       50,
		MaxDescriptionLength: 200,
		DateFormat:           time.RFC3339,
		AutoBackup:           true,
		BackupInterval:       24 * time.Hour,
	},
}

func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create default config if it doesn't exist
		return &cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
