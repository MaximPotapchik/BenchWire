package runner

import (
	"fmt"
	"math/rand/v2"
	"os/exec"
	"benchwire/internal/config"
)

type RunArgs struct {
	Methodology string
    Runs int
    CooldownMs int
    OutputDir string
    Targets []config.Target
}

// TODO: Make this work for more than just --benchmarks-file. Other benchmarkers.
func buildArgs(cfg config.Target, outputDir, prefix string, run int) []string {
	args := append([]string{}, cfg.Flags...)
	return append(args, fmt.Sprintf("--benchmarks-file=%s/%srun_%d.yaml", outputDir, prefix, run))
}

func execute(binary string, args []string) error {
	cmd := exec.Command(binary, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// TODO: - Default switch.
func Run(args RunArgs) error {

	switch args.Methodology {
	case "single":
		for i := 1; i <= args.Runs; i++ {
			if err := execute(args.Targets[0].BinPath, buildArgs(args.Targets[0], args.OutputDir, "", i)); err != nil {
				return err
			}
			fmt.Printf("Run %d/%d complete\n", i, args.Runs)
			msleep(args.CooldownMs)
		}

	case "sequential":
		for i := 1; i <= args.Runs; i++ {
			if err := execute(args.Targets[0].BinPath, buildArgs(args.Targets[0], args.OutputDir, "A", i)); err != nil {
				return err
			}
			msleep(args.CooldownMs)
		}
		for i := 1; i <= args.Runs; i++ {
			if err := execute(args.Targets[1].BinPath, buildArgs(args.Targets[1], args.OutputDir, "B", i)); err != nil {
				return err
			}
			msleep(args.CooldownMs)
		}

	case "cycling":
		for i := 1; i <= args.Runs; i++ {
			if err := execute(args.Targets[0].BinPath, buildArgs(args.Targets[0], args.OutputDir, "A", i)); err != nil {
				return err
			}
			msleep(args.CooldownMs)
			if err := execute(args.Targets[1].BinPath, buildArgs(args.Targets[1], args.OutputDir, "B", i)); err != nil {
				return err
			}
			msleep(args.CooldownMs)
		}

	case "random interleaving":
		order := make([]byte, 0, args.Runs*2)
		for i := 0; i < args.Runs; i++ {
			order = append(order, 'A', 'B')
		}
		rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

		countA, countB := 1, 1
		for _, side := range order {
			if side == 'A' {
				if err := execute(args.Targets[0].BinPath, buildArgs(args.Targets[0], args.OutputDir, "A", countA)); err != nil {
					return err
				}
				countA++
			} else {
				if err := execute(args.Targets[1].BinPath, buildArgs(args.Targets[1], args.OutputDir, "B", countB)); err != nil {
					return err
				}
				countB++
			}
			msleep(args.CooldownMs)
		}
	}
	return nil
}
