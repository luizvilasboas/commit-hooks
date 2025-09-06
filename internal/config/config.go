package config

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Types  []string `toml:"types"`
	Scopes []string `toml:"scopes"`
}

func Load() Config {
	defaultConfig := Config{
		Types:  []string{"feat", "fix", "docs", "style", "refactor", "test", "chore", "build"},
		Scopes: []string{"api", "ui", "db", "auth", "deps"},
	}

	localPath := filepath.Join("data", "conventional_commits.toml")
	if _, err := os.Stat(localPath); err == nil {
		var config Config
		if _, err := toml.DecodeFile(localPath, &config); err == nil {
			return config
		}
		log.Printf("Error decoding local config file %s: %v", localPath, err)
	}

	usr, err := user.Current()
	if err == nil {
		homePath := filepath.Join(usr.HomeDir, ".commit_hooks", "conventional_commits.toml")
		if _, err := os.Stat(homePath); err == nil {
			var config Config
			if _, err := toml.DecodeFile(homePath, &config); err == nil {
				return config
			}
			log.Printf("Error decoding global config file %s: %v", homePath, err)
		}
	}

	return defaultConfig
}
