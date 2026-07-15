package runner

import (
	"fmt"
	"math/rand/v2"
	"os/exec"
	"strings"
)

type config struct {
	binary string
	args []string
}

func buildConfigs(methodology string, args []string) (single, a, b config) {
	if strings.Contains(args[0], "|") {
		parts := strings.SplitN(args[0], "|", 2)
		a.binary, b.binary = parts[0], parts[1]
	} else {
		single.binary = args[0]
	}

	if methodology == "single" {
		single.args = append(single.args, args[1:]...)
		return
	}

	for _, opt := range args[1:] {
		lhs, val, found := strings.Cut(opt, "=")
		if !found || lhs == "" {
			continue
		}
		marker := lhs[len(lhs)-1]
		flagName := lhs[:len(lhs)-1]
		if marker == 'A' {
			a.args = append(a.args, flagName+"="+val)
		} else {
			b.args = append(b.args, flagName+"="+val)
		}
	}
	return
}

func buildArgs(cfg config, outputDir, prefix string, run int) []string {
	args := append([]string{}, cfg.args...)
	return append(args, fmt.Sprintf("--benchmarks-file=%s/%srun_%d.yaml", outputDir, prefix, run))
}

func execute(binary string, args []string) error {
	cmd := exec.Command(binary, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// TODO: - Default switch.
func Run(methodology string, runs, cooldownMs int, outputDir string, args []string) error {
	single, a, b := buildConfigs(methodology, args)

	switch methodology {
	case "single":
		for i := 1; i <= runs; i++ {
			if err := execute(single.binary, buildArgs(single, outputDir, "", i)); err != nil {
				return err
			}
			fmt.Printf("Run %d/%d complete\n", i, runs)
			msleep(cooldownMs)
		}

	case "sequential":
		for i := 1; i <= runs; i++ {
			if err := execute(a.binary, buildArgs(a, outputDir, "A", i)); err != nil {
				return err
			}
			msleep(cooldownMs)
		}
		for i := 1; i <= runs; i++ {
			if err := execute(b.binary, buildArgs(b, outputDir, "B", i)); err != nil {
				return err
			}
			msleep(cooldownMs)
		}

	case "cycling":
		for i := 1; i <= runs; i++ {
			if err := execute(a.binary, buildArgs(a, outputDir, "A", i)); err != nil {
				return err
			}
			msleep(cooldownMs)
			if err := execute(b.binary, buildArgs(b, outputDir, "B", i)); err != nil {
				return err
			}
			msleep(cooldownMs)
		}

	case "random interleaving":
		order := make([]byte, 0, runs*2)
		for i := 0; i < runs; i++ {
			order = append(order, 'A', 'B')
		}
		rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

		countA, countB := 1, 1
		for _, side := range order {
			if side == 'A' {
				if err := execute(a.binary, buildArgs(a, outputDir, "A", countA)); err != nil {
					return err
				}
				countA++
			} else {
				if err := execute(b.binary, buildArgs(b, outputDir, "B", countB)); err != nil {
					return err
				}
				countB++
			}
			msleep(cooldownMs)
		}
	}
	return nil
}
