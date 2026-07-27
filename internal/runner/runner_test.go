package runner

import (
	"reflect"
	"testing"

	"benchwire/internal/config"
)

func TestBuildArgsSingleMode(t *testing.T) {
	cfg := config.Target{Label: "Raw", BinPath: "/bin/exegesis", Flags: []string{"--mcpu=native", "--mode=latency"}}
	got := buildArgs(cfg, "/out", "", 3)
	want := []string{"--mcpu=native", "--mode=latency", "--benchmarks-file=/out/run_3.yaml"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBuildArgsComparePrefix(t *testing.T) {
	cases := []struct {
		name   string
		prefix string
		run    int
		want   string
	}{
		{"side A", "A", 1, "--benchmarks-file=/out/Arun_1.yaml"},
		{"side B", "B", 5, "--benchmarks-file=/out/Brun_5.yaml"},
	}

	cfg := config.Target{Flags: []string{"--mcpu=native"}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := buildArgs(cfg, "/out", c.prefix, c.run)
			last := got[len(got)-1]
			if last != c.want {
				t.Errorf("got %q, want %q", last, c.want)
			}
		})
	}
}

func TestBuildArgsDoesNotMutateOriginalFlags(t *testing.T) {
	cfg := config.Target{Flags: []string{"--mcpu=native"}}
	buildArgs(cfg, "/out", "", 1)
	buildArgs(cfg, "/out", "", 2)

	if len(cfg.Flags) != 1 {
		t.Errorf("cfg.Flags got mutated across calls, now has %d entries, want 1", len(cfg.Flags))
	}
}
