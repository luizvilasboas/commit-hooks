package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func createTempConfigFile(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "conventional_commits.toml")
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		t.Fatalf("Failed to create parent directories for temp config: %v", err)
	}
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}
	return path
}

func TestLoad(t *testing.T) {
	defaultConfig := Config{
		Types:  []string{"feat", "fix", "docs", "style", "refactor", "test", "chore", "build"},
		Scopes: []string{"api", "ui", "db", "auth", "deps"},
	}

	t.Run("should return default config when no files exist", func(t *testing.T) {
		tempDir := t.TempDir()
		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)

		t.Setenv("HOME", tempDir)

		cfg := Load()

		if !reflect.DeepEqual(cfg, defaultConfig) {
			t.Errorf("Expected default config, got %+v", cfg)
		}
	})

	t.Run("should load local config file when it exists", func(t *testing.T) {
		tempDir := t.TempDir()
		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)
		t.Setenv("HOME", t.TempDir())

		localConfigContent := `
			types = ["local_feat", "local_fix"]
			scopes = ["local_api", "local_ui"]
		`

		createTempConfigFile(t, filepath.Join(tempDir, "data"), localConfigContent)

		expectedConfig := Config{
			Types:  []string{"local_feat", "local_fix"},
			Scopes: []string{"local_api", "local_ui"},
		}

		cfg := Load()

		if !reflect.DeepEqual(cfg, expectedConfig) {
			t.Errorf("Expected local config %+v, got %+v", expectedConfig, cfg)
		}
	})

	t.Run("should load global config file when local does not exist", func(t *testing.T) {
		projectDir := t.TempDir()
		originalWd, _ := os.Getwd()
		os.Chdir(projectDir)
		defer os.Chdir(originalWd)

		tempHomeDir := t.TempDir()
		t.Setenv("HOME", tempHomeDir)

		globalConfigContent := `
			types = ["global_feat", "global_fix"]
			scopes = ["global_api", "global_ui"]
		`

		createTempConfigFile(t, filepath.Join(tempHomeDir, ".commit_hooks"), globalConfigContent)

		expectedConfig := Config{
			Types:  []string{"global_feat", "global_fix"},
			Scopes: []string{"global_api", "global_ui"},
		}

		cfg := Load()

		if !reflect.DeepEqual(cfg, expectedConfig) {
			t.Errorf("Expected global config %+v, got %+v", expectedConfig, cfg)
		}
	})

	t.Run("should prioritize local config over global config", func(t *testing.T) {
		projectDir := t.TempDir()
		originalWd, _ := os.Getwd()
		os.Chdir(projectDir)
		defer os.Chdir(originalWd)

		tempHomeDir := t.TempDir()
		t.Setenv("HOME", tempHomeDir)
		globalConfigContent := `
			types = ["global_feat"]
			scopes = ["global_api"]
		`

		createTempConfigFile(t, filepath.Join(tempHomeDir, ".commit_hooks"), globalConfigContent)

		localConfigContent := `
			types = ["local_feat"]
			scopes = ["local_api"]
		`

		createTempConfigFile(t, filepath.Join(projectDir, "data"), localConfigContent)

		expectedConfig := Config{
			Types:  []string{"local_feat"},
			Scopes: []string{"local_api"},
		}

		cfg := Load()

		if !reflect.DeepEqual(cfg, expectedConfig) {
			t.Errorf("Expected local config to be prioritized, got %+v", cfg)
		}
	})

	t.Run("should fall back to default when local config is malformed", func(t *testing.T) {
		tempDir := t.TempDir()
		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)
		t.Setenv("HOME", t.TempDir())

		malformedContent := `types = ["invalid_toml`
		createTempConfigFile(t, filepath.Join(tempDir, "data"), malformedContent)

		cfg := Load()

		if !reflect.DeepEqual(cfg, defaultConfig) {
			t.Errorf("Expected default config after malformed local file, got %+v", cfg)
		}
	})
}
