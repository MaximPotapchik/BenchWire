package runner

import (
	"benchwire/internal/config"
	"fmt"
	"math/rand/v2"
	"os/exec"
)

type RunArgs struct {
	Methodology string
    Runs int
    Cooldown config.CooldownTimer
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

	totalRuns := args.Runs
	if args.Methodology != "single" {
		totalRuns = args.Runs * 2
	}

	progress := ProgressBar {
		TotalRuns: totalRuns,
		RunNumber: 0,
		CapturedTimes: []RunTiming{},
		RunningSum: 0,
		CooldownSum: 0,
		EstimatedETA: 0,
		Percent: 0,
	}

	switch args.Methodology {
		case "single":
			for i := 1; i <= args.Runs; i++ {
				timerStart := GetTime()
				if err := execute(args.Targets[0].BinPath, buildArgs(args.Targets[0], args.OutputDir, "", i)); err != nil {
					return err
				}
				timerEnd := GetTime()
				cooldown, err := msleep(args.Cooldown)
				if err != nil {
					return err
				}
				progress.Tick(i, timerStart, timerEnd, cooldown)
			}

		case "sequential":
			for i := 1; i <= args.Runs; i++ {
				timerStart := GetTime()
				if err := execute(args.Targets[0].BinPath, buildArgs(args.Targets[0], args.OutputDir, "A", i)); err != nil {
					return err
				}
				timerEnd := GetTime()
				cooldown, err := msleep(args.Cooldown)
				if err != nil {
					return err
				}
				progress.Tick(i, timerStart, timerEnd, cooldown)
			}
			for i := 1; i <= args.Runs; i++ {
				timerStart := GetTime()
				if err := execute(args.Targets[1].BinPath, buildArgs(args.Targets[1], args.OutputDir, "B", i)); err != nil {
					return err
				}
				timerEnd := GetTime()
				cooldown, err := msleep(args.Cooldown)
				if err != nil {
					return err
				}
				progress.Tick(i, timerStart, timerEnd, cooldown)
			}

		case "cycling":
			for i := 1; i <= args.Runs; i++ {
				timerStart := GetTime()
				if err := execute(args.Targets[0].BinPath, buildArgs(args.Targets[0], args.OutputDir, "A", i)); err != nil {
					return err
				}
				timerEnd := GetTime()
				cooldown, err := msleep(args.Cooldown)
				if err != nil {
					return err
				}
				progress.Tick(i, timerStart, timerEnd, cooldown)

				timerStart = GetTime()
				if err := execute(args.Targets[1].BinPath, buildArgs(args.Targets[1], args.OutputDir, "B", i)); err != nil {
					return err
				}
				timerEnd = GetTime()
				cooldown, err = msleep(args.Cooldown)
				if err != nil {
					return err
				}
				progress.Tick(i, timerStart, timerEnd, cooldown)
			}

		case "random interleaving":
			order := make([]byte, 0, args.Runs*2)
			for i := 0; i < args.Runs; i++ {
				order = append(order, 'A', 'B')
			}
			rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

			countA, countB := 1, 1
			for idx, side := range order {
				timerStart := GetTime()
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
				timerEnd := GetTime()
				cooldown, err := msleep(args.Cooldown)
				if err != nil {
					return err
				}
				progress.Tick(idx+1, timerStart, timerEnd, cooldown)
			}
	}
	fmt.Printf("\n[BenchWire] total: %.2fs\n", float64(GetTotalTimeSpent(progress))/1e9)
	return nil
}
