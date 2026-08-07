package runner

import (
	"benchwire/internal/scheduler"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
)

// TODO: Make this work for more than just --benchmarks-file. Other benchmarkers.
// Meaning, it needs to become preset based.
func buildArgs(flags []string, outputDir string, matrixName string, label string, prefix string, run int) []string {
	scheduled := append([]string{}, flags...)
	return append(scheduled, fmt.Sprintf("--benchmarks-file=%s/%s_%s_%srun%d.yaml", outputDir, 
									     matrixName, label,  prefix, run))
}

func execute(binary string, scheduled []string) error {
	cmd := exec.Command(binary, scheduled...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TODO: - Default switch.
// - Create a function for each loop.
func Run(scheduled []scheduler.ScheduledMatrix, outputDir string) error {
	
	var totalRuns int
	for _, specMatrix := range scheduled {
		for _, targets := range specMatrix.RunTargets {
			if targets.Methodology == "single" {
				totalRuns += targets.Runs
			} else {
				totalRuns += targets.Runs * 2
			}
		}
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
	
	runNumber := 0

	for _, specMatrix := range scheduled{
		for _, runTarget := range specMatrix.RunTargets{
			switch runTarget.Methodology {
				case "single":
					for i := 1; i <= runTarget.Runs; i++ {
						timerStart := GetTime()
						if err := execute(runTarget.BinPath[0], buildArgs(runTarget.Flags[0], outputDir, specMatrix.SpecMatrixName, runTarget.Label[0], "", i)); err != nil {
							return err
						}
						timerEnd := GetTime()

						cooldown, err := msleep(runTarget.CooldownTimer)
						if err != nil {
							return err
						}

						runNumber++
						progress.Tick(runNumber, timerStart, timerEnd, cooldown, specMatrix.SpecMatrixName)
					}

				case "sequential":
					for i := 1; i <= runTarget.Runs; i++ {
						timerStart := GetTime()
						if err := execute(runTarget.BinPath[0], buildArgs(runTarget.Flags[0], outputDir, specMatrix.SpecMatrixName, runTarget.Label[0], "A", i)); err != nil {
							return err
						}
						timerEnd := GetTime()

						cooldown, err := msleep(runTarget.CooldownTimer)
						if err != nil {
							return err
						}

						runNumber++
						progress.Tick(runNumber, timerStart, timerEnd, cooldown, specMatrix.SpecMatrixName)
					}
					for i := 1; i <= runTarget.Runs; i++ {
						timerStart := GetTime()
						if err := execute(runTarget.BinPath[1], buildArgs(runTarget.Flags[1], outputDir, specMatrix.SpecMatrixName, runTarget.Label[1], "B", i)); err != nil {
							return err
						}
						timerEnd := GetTime()

						cooldown, err := msleep(runTarget.CooldownTimer)
						if err != nil {
							return err
						}

						runNumber++
						progress.Tick(runNumber, timerStart, timerEnd, cooldown, specMatrix.SpecMatrixName)
					}

				case "cycling":
					for i := 1; i <= runTarget.Runs; i++ {
						timerStart := GetTime()
						if err := execute(runTarget.BinPath[0], buildArgs(runTarget.Flags[0], outputDir, specMatrix.SpecMatrixName, runTarget.Label[0], "A", i)); err != nil {
							return err
						}
						timerEnd := GetTime()

						cooldown, err := msleep(runTarget.CooldownTimer)
						if err != nil {
							return err
						}

						runNumber++
						progress.Tick(runNumber, timerStart, timerEnd, cooldown, specMatrix.SpecMatrixName)

						timerStart = GetTime()
						if err := execute(runTarget.BinPath[1], buildArgs(runTarget.Flags[1], outputDir, specMatrix.SpecMatrixName, runTarget.Label[1], "B", i)); err != nil {
							return err
						}
						timerEnd = GetTime()

						cooldown, err = msleep(runTarget.CooldownTimer)
						if err != nil {
							return err
						}

						runNumber++
						progress.Tick(runNumber, timerStart, timerEnd, cooldown, specMatrix.SpecMatrixName)
					}

				case "random interleaving":
					order := make([]byte, 0, runTarget.Runs*2)
					for i := 0; i < runTarget.Runs; i++ {
						order = append(order, 'A', 'B')
					}
					rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

					countA, countB := 1, 1
					for _, side := range order {
						timerStart := GetTime()
						if side == 'A' {
							if err := execute(runTarget.BinPath[0], buildArgs(runTarget.Flags[0], outputDir, specMatrix.SpecMatrixName, runTarget.Label[0], "A", countA)); err != nil {
								return err
							}
							countA++
						} else {
							if err := execute(runTarget.BinPath[1], buildArgs(runTarget.Flags[1], outputDir, specMatrix.SpecMatrixName, runTarget.Label[1], "B", countB)); err != nil {
								return err
							}
							countB++
						}
						timerEnd := GetTime()

						cooldown, err := msleep(runTarget.CooldownTimer)
						if err != nil {
							return err
						}

						runNumber++
						progress.Tick(runNumber, timerStart, timerEnd, cooldown, specMatrix.SpecMatrixName)
					}
				}
			}
	}
	fmt.Printf("\n[BenchWire] total: %.2fs\n", float64(GetTotalTimeSpent(progress))/1e9)
	return nil
}
