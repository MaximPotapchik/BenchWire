package scheduler

import (
	"benchwire/internal/config"
	"benchwire/internal/resolver"
	"testing"
)

func TestScheduler(t *testing.T) {
	cfg, err := config.LoadYamlConfig("testdata")

	if err != nil {
		t.Fatalf("Expected to find and parse the file, got error: %v", err)
	}

	resolved := resolver.Resolve(*cfg)
	scheduled, err := Schedule(resolved)

	if err != nil {
		t.Fatalf("Expected Schedule to succeed, got error: %v", err)
	}

	if len(scheduled) != 2 {
		t.Fatalf("Expected 2 scheduled matrices. Got: %d", len(scheduled))
	}

	expectedMatrices := []struct {
		Name string
		Methodology string
		Runs int
		BinPaths map[string]string
	}{
		{
			Name: "VEXTRACTF128rri FPU pipe quick check",
			Methodology: "sequential",
			Runs: 20,
			BinPaths: map[string]string{ "build1": "TestBin1", "build2": "TestBin2" },
		},
		{
			Name: "VEXTRACTF128rri FPU pipe full validation",
			Methodology: "random interleaving",
			Runs: 100,
			BinPaths: map[string]string{ "build1": "TestBin1", "build2": "TestBin2" },
		},
	}

	for i, expected := range expectedMatrices {
		if scheduled[i].SpecMatrixName != expected.Name {
			t.Fatalf("Expected matrix %d Name to be %s. Got: %s", i, expected.Name, scheduled[i].SpecMatrixName)
		}

		if scheduled[i].Benchmarker != "llvm-exegesis" {
			t.Fatalf("Expected matrix %d Benchmarker to be llvm-exegesis. Got: %s", i, scheduled[i].Benchmarker)
		}

		if len(scheduled[i].RunTargets) != 3 {
			t.Fatalf("Expected matrix %d to have 3 RunTargets (3 sequence steps). Got: %d", i, len(scheduled[i].RunTargets))
		}

		for j, rt := range scheduled[i].RunTargets {
			if rt.Runs != expected.Runs {
				t.Fatalf("Matrix %d RunTarget %d: expected Runs %d. Got: %d", i, j, expected.Runs, rt.Runs)
			}

			if len(rt.Label) != len(rt.BinPath) || len(rt.Label) != len(rt.Flags) {
				t.Fatalf("Matrix %d RunTarget %d: Label/BinPath/Flags length mismatch. Label:%d BinPath:%d Flags:%d", i, j, len(rt.Label), len(rt.BinPath), len(rt.Flags))
			}

			isSolo := len(rt.Label) == 1

			if isSolo && rt.Methodology != "single" {
				t.Fatalf("Matrix %d RunTarget %d: solo step should force Methodology single. Got: %s", i, j, rt.Methodology)
			}

			if !isSolo && rt.Methodology != expected.Methodology {
				t.Fatalf("Matrix %d RunTarget %d: pair step should keep matrix Methodology %s. Got: %s", i, j, expected.Methodology, rt.Methodology)
			}

			for k, label := range rt.Label {
				wantBinPath, ok := expected.BinPaths[label]
				if !ok {
					t.Fatalf("Matrix %d RunTarget %d: unexpected label %q not in testdata", i, j, label)
				}
				if rt.BinPath[k] != wantBinPath {
					t.Fatalf("Matrix %d RunTarget %d: label %q expected BinPath %s. Got: %s", i, j, label, wantBinPath, rt.BinPath[k])
				}
			}
		}
	}
}
