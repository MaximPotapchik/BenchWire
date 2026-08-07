package config

import (
	"testing"
)

func TestLoadYaml(t *testing.T) {
	cfg, err := LoadYamlConfig("testdata")

	if err != nil {
		t.Fatalf("Expected to find and parse the file, got error: %v", err)
	}

	if cfg.Default.Methodology != "single" {
		t.Fatalf("Expected Methodology to be single. Got: %s", cfg.Default.Methodology)
	}

	if cfg.Default.Runs != 50 {
		t.Fatalf("Expected CooldownTimer to be 50. Got: %d", cfg.Default.Runs)
	}

	if cfg.Default.CooldownTimer.Value != "5ms" {
		t.Fatalf("Expected Methodology to be 5ms. Got: %s", cfg.Default.CooldownTimer.Value)
	}

	if cfg.GlobalFlags[0] != "--mcpu=native" {
		t.Fatalf("Expected GlobalFlags[0] to be --mcpu=native. Got: %s", cfg.GlobalFlags[0])
	}
	
	//specMatrix
	if cfg.SpecMatrix[0].Name != "VEXTRACTF128rri FPU pipe quick check" {
		t.Fatalf("Expected name to be VEXTRACTF128rri. Got: %s", cfg.SpecMatrix[0].Name)
	}

	if cfg.SpecMatrix[0].Benchmarker != "llvm-exegesis" {
		t.Fatalf("Expected benchmarker to be llvm-exegesis. Got: %s", cfg.SpecMatrix[0].Benchmarker)
	}
	
	if cfg.SpecMatrix[0].Methodology != "sequential" {
		t.Fatalf("Expected Methodology to be sequential. Got: %s", cfg.SpecMatrix[0].Methodology)
	}
	
	if cfg.SpecMatrix[0].Runs != 20 {
		t.Fatalf("Expected Mruns to be 20. Got: %d", cfg.SpecMatrix[0].Runs)
	}

	if len(cfg.SpecMatrix[0].CooldownTimer.Value) != 0 {
		t.Fatalf("Expected Methodology to be VEXTRACTF128rri. Got: %s", cfg.SpecMatrix[0].Name)
	}
	
	if cfg.SpecMatrix[0].LocalFlags[0] != "--mode=uops" {
		t.Fatalf("Expected Methodology to be --mode=uops. Got: %s", cfg.SpecMatrix[0].LocalFlags[0]) 
	}

	// SpecMatrix Presets
	if cfg.SpecMatrix[0].Presets[0].Name != "VexLowMin" {
		t.Fatalf("Expected Methodology to be VexLowMin. Got: %s", cfg.SpecMatrix[0].Presets[0].Name)
	}

	if cfg.SpecMatrix[0].Presets[0].Inherit[0] != "Vex" {
		t.Fatalf("Expected Methodology to be Vex. Got: %s",  cfg.SpecMatrix[0].Presets[0].Inherit[0])
	}

	if cfg.SpecMatrix[0].Presets[0].Flags[0] != "--min-instructions=5000" {
		t.Fatalf("Expected Flags to be --min... . Got: %v", cfg.SpecMatrix[0].Presets[0].Flags[0])
	}

	// SpecMatrix Targets
	if cfg.SpecMatrix[0].Targets[0].Label != "build1" {
		t.Fatalf("Expected Label to be Megapatch. Got: %s", cfg.SpecMatrix[0].Targets[0].Label)
	}

	if cfg.SpecMatrix[0].Targets[0].BinPath != "TestBin1" {
		t.Fatalf("Expected BinPath to be TestBin1. Got: %s", cfg.SpecMatrix[0].Targets[0].BinPath)
	}

	if cfg.SpecMatrix[0].Targets[0].Preset[0] != "VexLowMin" {
		t.Fatalf("Expected Preset[0] to be VexLowMin. Got: %s", cfg.SpecMatrix[0].Targets[0].Preset[0])
	}

	if len(cfg.SpecMatrix[0].Targets[0].Flags) != 0 {
		t.Fatalf("Expected Flags to be empty. Got: %v", cfg.SpecMatrix[0].Targets[0].Flags)
	}

	if cfg.SpecMatrix[0].Targets[1].Label != "build2" {
		t.Fatalf("Expected Label to be libpfm. Got: %s", cfg.SpecMatrix[0].Targets[1].Label)
	}

	if cfg.SpecMatrix[0].Targets[1].BinPath != "TestBin2" {
		t.Fatalf("Expected BinPath to be TestBin2. Got: %s", cfg.SpecMatrix[0].Targets[1].BinPath)
	}

	if cfg.SpecMatrix[0].Targets[1].Preset[0] != "VexLowMin" {
		t.Fatalf("Expected Preset[0] to be VexLowMin. Got: %s", cfg.SpecMatrix[0].Targets[1].Preset[0])
	}

	if cfg.SpecMatrix[0].Targets[1].Flags[0] != "--validation-counter=l1d-cache" {
		t.Fatalf("Expected Flags[0] to be --validation-counter=l1d-cache. Got: %s", cfg.SpecMatrix[0].Targets[1].Flags[0])
	}

	// SpecMatrix Sequence
	if len(cfg.SpecMatrix[0].Sequence) != 3 {
		t.Fatalf("Expected Sequence to have 3 steps. Got: %d", len(cfg.SpecMatrix[0].Sequence))
	}

	if len(cfg.SpecMatrix[0].Sequence[0]) != 1 || cfg.SpecMatrix[0].Sequence[0][0] != "build1" {
		t.Fatalf("Expected Sequence[0] to be [Megapatch]. Got: %v", cfg.SpecMatrix[0].Sequence[0])
	}

	if len(cfg.SpecMatrix[0].Sequence[1]) != 1 || cfg.SpecMatrix[0].Sequence[1][0] != "build2" {
		t.Fatalf("Expected Sequence[1] to be [libpfm]. Got: %v", cfg.SpecMatrix[0].Sequence[1])
	}

	if len(cfg.SpecMatrix[0].Sequence[2]) != 2 || cfg.SpecMatrix[0].Sequence[2][0] != "build1" || cfg.SpecMatrix[0].Sequence[2][1] != "build2" {
		t.Fatalf("Expected Sequence[2] to be [build1 build2]. Got: %v", cfg.SpecMatrix[0].Sequence[2])
	}

	if cfg.SpecMatrix[1].Name != "VEXTRACTF128rri FPU pipe full validation" {
		t.Fatalf("Expected specMatrix to be VEXTRACTF128rri. Got: %s", cfg.SpecMatrix[1].Name )
	}

}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := LoadYamlConfig("testdata/nonexistent_dir")

	if err == nil {
		t.Fatal("Expected an error for a missing config.yaml, got nil")
	}
}

