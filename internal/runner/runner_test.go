package runner

import (
	"reflect"
	"testing"
)

//TODO: Test Run() itself. 

func TestBuildArgsSingleMode(t *testing.T) {
	flags := []string{"--mcpu=native", "--mode=latency"}
	got := buildArgs(flags, "/out", "MatrixA", "build1", "", 3)
	want := []string{"--mcpu=native", "--mode=latency", "--benchmarks-file=/out/MatrixA_build1_run3.yaml"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgsComparePrefix(t *testing.T) {
	cases := []struct {
		name string
		matrixName string
		label string
		prefix string
		run int
		want string
	}{
		{"side A", "MatrixA", "build1", "A", 1, "--benchmarks-file=/out/MatrixA_build1_Arun1.yaml"},
		{"side B", "MatrixB", "build2", "B", 5, "--benchmarks-file=/out/MatrixB_build2_Brun5.yaml"},
	}

	flags := []string{"--mcpu=native"}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := buildArgs(flags, "/out", c.matrixName, c.label, c.prefix, c.run)
			last := got[len(got)-1]
			if last != c.want {
				t.Errorf("got %q, want %q", last, c.want)
			}
		})
	}
}

func TestBuildArgsDoesNotMutateOriginalFlags(t *testing.T) {
	flags := []string{"--mcpu=native"}
	buildArgs(flags, "/out", "MatrixA", "build1", "", 1)
	buildArgs(flags, "/out", "MatrixA", "build2", "", 2)

	if len(flags) != 1 {
		t.Errorf("flags got mutated across calls, now has %d entries, want 1", len(flags))
	}
}
