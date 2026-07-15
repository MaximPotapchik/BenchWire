package runner

import "testing"

func TestBuildConfigsSingle(t *testing.T) {
	args := []string{"/path/to/exegesis", "--mode=latency"}
	single, _, _ := buildConfigs("single", args)

	if single.binary != "/path/to/exegesis" {
		t.Errorf("expected %q, got %q", "/path/to/exegesis", single.binary)
	}
}

func TestBuildConfigsCompare(t *testing.T) {
	cases := []struct {
		name string
		args []string
		wantABin string
		wantBBin string
	}{
		{
			name: "basic A/B split",
			args: []string{"binA|binB", "--mcpuA=native", "--mcpuB=native"},
			wantABin: "binA",
			wantBBin: "binB",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, a, b := buildConfigs("sequential", c.args)
			if a.binary != c.wantABin {
				t.Errorf("A: want %q, got %q", c.wantABin, a.binary)
			}
			if b.binary != c.wantBBin {
				t.Errorf("B: want %q, got %q", c.wantBBin, b.binary)
			}
		})
	}
}
