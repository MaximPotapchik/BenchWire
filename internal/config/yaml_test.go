package config

import (
	"testing"
)

func TestLoadConfigFindsFile(t *testing.T) {
	cfg, err := LoadYamlConfig("testdata")

	if err != nil {
		t.Fatalf("expected to find and parse the file, got error: %v", err)
	}

	if cfg.Methodology != "random interleaving" {
		t.Errorf("expected methodology 'random interleaving', got %q", cfg.Methodology)
	}

	if cfg.Runs != 5 {
		t.Errorf("expected Runs=5, got %d", cfg.Runs)
	}

	if len(cfg.Targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(cfg.Targets))
	}

	if cfg.Targets[0].Label != "Raw" {
		t.Errorf("expected label 'Raw', got %q", cfg.Targets[0].Label)
	}

	if cfg.Targets[0].BinPath != "/bin/raw" {
		t.Errorf("expected BinPath '/bin/raw', got %q (empty means the yaml key doesn't match the struct tag)", cfg.Targets[0].BinPath)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := LoadYamlConfig("testdata/nonexistent_dir")

	if err == nil {
		t.Fatal("expected an error for a missing config.yaml, got nil")
	}
}

