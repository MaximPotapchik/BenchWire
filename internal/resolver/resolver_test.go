package resolver

import (
	"testing"
	"benchwire/internal/config"
)

func TestResolver(t *testing.T) {
	cfg, err := config.LoadYamlConfig("testdata")

	if err != nil {
		t.Fatalf("Expected to find and parse the file, got error: %v", err)
	}

	resolved := Resolve(*cfg)
	
	if len(resolved) == 0 {
		t.Fatalf("Expected at least one resolved specMatrix, got 0")
	}

	if resolved[0].Name != "VEXTRACTF128rri FPU pipe quick check" {
		t.Fatalf("Expected Name to be VEXTRACTF128rri FPU pipe quick check. Got: %s", resolved[0].Name)
	}

	if resolved[0].Benchmarker != "llvm-exegesis" {
		t.Fatalf("Expected Benchmarker to be llvm-exegesis. Got: %s", resolved[0].Benchmarker)
	}

	if resolved[0].Methodology != "sequential" {
		t.Fatalf("Expected Methodology to be sequential. Got: %s", resolved[0].Methodology)
	}

	if resolved[0].Runs != 20 {
		t.Fatalf("Expected Runs to be 20. Got: %d", resolved[0].Runs)
	}

	if resolved[0].CooldownTimer.Value != "5ms" {
		t.Fatalf("Expected CooldownTimer to fall back to default 5ms. Got: %s", resolved[0].CooldownTimer.Value)
	}

	if len(resolved[0].ResolvedTarget) != 2 {
		t.Fatalf("Expected 2 resolved targets. Got: %d", len(resolved[0].ResolvedTarget))
	}

	if resolved[0].ResolvedTarget[0].Label != "build1" {
		t.Fatalf("Expected Label to be build1. Got: %s", resolved[0].ResolvedTarget[0].Label)
	}

	if resolved[0].ResolvedTarget[0].BinPath != "TestBin1" {
		t.Fatalf("Expected BinPath to be TestBin1. Got: %s", resolved[0].ResolvedTarget[0].BinPath)
	}

	expectedbuild1 := []string{"--mcpu=native", "--mode=uops", "--opcode-name=VEXTRACTF128rri", "--min-instructions=5000"}
	if len(resolved[0].ResolvedTarget[0].Flags) != len(expectedbuild1) {
		t.Fatalf("Expected build1 Flags %v. Got: %v", expectedbuild1, resolved[0].ResolvedTarget[0].Flags)
	}
	for i, f := range expectedbuild1 {
		if resolved[0].ResolvedTarget[0].Flags[i] != f {
			t.Fatalf("Expected build1 Flags[%d] to be %s. Got: %s", i, f, resolved[0].ResolvedTarget[0].Flags[i])
		}
	}

	if resolved[0].ResolvedTarget[1].Label != "build2" {
		t.Fatalf("Expected Label to be build2. Got: %s", resolved[0].ResolvedTarget[1].Label)
	}

	expectedbuild2 := []string{"--mcpu=native", "--mode=uops", "--opcode-name=VEXTRACTF128rri", "--min-instructions=5000", "--validation-counter=l1d-cache"}
	if len(resolved[0].ResolvedTarget[1].Flags) != len(expectedbuild2) {
		t.Fatalf("Expected build2 Flags %v. Got: %v", expectedbuild2, resolved[0].ResolvedTarget[1].Flags)
	}
	for i, f := range expectedbuild2{
		if resolved[0].ResolvedTarget[1].Flags[i] != f {
			t.Fatalf("Expected build2 Flags[%d] to be %s. Got: %s", i, f, resolved[0].ResolvedTarget[1].Flags[i])
		}
	}

	if len(resolved[0].Sequence) != 3 {
		t.Fatalf("Expected Sequence to have 3 steps. Got: %d", len(resolved[0].Sequence))
	}
}
