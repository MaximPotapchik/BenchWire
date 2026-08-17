package runner

import (
	"reflect"
	"testing"
)

//TODO: Test Run() itself.

func TestBuildArgsSingleMode(t *testing.T) {
	flags := []string{"--mcpu=native", "--mode=latency"}
	preset := outputPresets["llvm-exegesis"]
	got, path := buildArgs(flags, preset, "/out", "MatrixA", "build1", "", 3)
	want := []string{"--mcpu=native", "--mode=latency", "--benchmarks-file=/out/MatrixA_build1_run3.yaml"}
	wantPath := "/out/MatrixA_build1_run3.yaml"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if path != wantPath {
		t.Errorf("got path %q, want %q", path, wantPath)
	}
}

func TestBuildArgsComparePrefix(t *testing.T) {
	cases := []struct {
		name       string
		matrixName string
		label      string
		prefix     string
		run        int
		want       string
	}{
		{"side A", "MatrixA", "build1", "A", 1, "--benchmarks-file=/out/MatrixA_build1_Arun1.yaml"},
		{"side B", "MatrixB", "build2", "B", 5, "--benchmarks-file=/out/MatrixB_build2_Brun5.yaml"},
	}

	flags := []string{"--mcpu=native"}
	preset := outputPresets["llvm-exegesis"]

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := buildArgs(flags, preset, "/out", c.matrixName, c.label, c.prefix, c.run)
			last := got[len(got)-1]
			if last != c.want {
				t.Errorf("got %q, want %q", last, c.want)
			}
		})
	}
}

func TestBuildArgsDoesNotMutateOriginalFlags(t *testing.T) {
	flags := []string{"--mcpu=native"}
	preset := outputPresets["llvm-exegesis"]
	buildArgs(flags, preset, "/out", "MatrixA", "build1", "", 1)
	buildArgs(flags, preset, "/out", "MatrixA", "build2", "", 2)

	if len(flags) != 1 {
		t.Errorf("flags got mutated across calls, now has %d entries, want 1", len(flags))
	}
}

func TestBuildArgsPyperfPreset(t *testing.T) {
	flags := []string{"--fast"}
	preset := outputPresets["pyperf"]
	got, path := buildArgs(flags, preset, "/out", "MatrixA", "build1", "", 1)
	wantLast := "-o=/out/MatrixA_build1_run1.json"
	wantPath := "/out/MatrixA_build1_run1.json"

	last := got[len(got)-1]
	if last != wantLast {
		t.Errorf("got %q, want %q", last, wantLast)
	}
	if path != wantPath {
		t.Errorf("got path %q, want %q", path, wantPath)
	}
}
