package runner

import (
	"benchwire/internal/config"
	"benchwire/internal/scheduler"
	"bytes"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"strings"
	"path/filepath"
)

type OutputPreset struct {
	Execution string
	BuildArg func(outputDir, matrixName, label, prefix string, run int) string
}

var outputPresets = map[string]OutputPreset{
	"llvm-exegesis": {
		Execution: "",
		BuildArg: func(outputDir, matrixName, label, prefix string, run int) string {
			return fmt.Sprintf("--benchmarks-file=%s/%s_%s_%srun%d.yaml", outputDir, matrixName, label, prefix, run)
		},
	},
	"pyperf": {
		Execution: "python3",
		BuildArg: func(outputDir, matrixName, label, prefix string, run int) string {
			return fmt.Sprintf("-o=%s/%s_%s_%srun%d.json", outputDir, matrixName, label, prefix, run)
		},
	},
}

func buildArgs(flags []string, preset OutputPreset, outputDir, matrixName, label, prefix string, run int) ([]string, string) {
	flagString := preset.BuildArg(outputDir, matrixName, label, prefix, run)
	scheduled := append([]string{}, flags...)
	scheduled = append(scheduled, flagString)
	_, path, _ := strings.Cut(flagString, "=")

	return scheduled, path
}

func execute(execution, binary string, args []string) (string, error) {
	command := binary
	cmdArgs := args
	if execution != "" {
		command = execution
		cmdArgs = append([]string{binary}, args...)
	}

	cmd := exec.Command(command, cmdArgs...)
	var stderrBuf bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderrBuf
	err := cmd.Run()

	return stderrBuf.String(), err
}

// Make one per error. 
func WriteFailureYaml(outputPath string, errText string, opcode string) error {
	content := fmt.Sprintf("error: %q\nmeasurements: []\nkey:\n  instructions:\n    - %q\n", errText, opcode)
	return os.WriteFile(outputPath, []byte(content), 0644)
}

func SingularRun(progress *ProgressBar, runNumber int, outputPreset OutputPreset, outputDir, specMatrixName, binPath string, flags []string, label, prefix string, run int, cooldownTimer config.CooldownTimer) (int, error) {
	timerStart := GetTime()
	args, outputPath := buildArgs(flags, outputPreset, outputDir, specMatrixName, label, prefix, run)
	os.Remove(outputPath)
	stderrText, execErr := execute(outputPreset.Execution, binPath, args)

	if _, statErr := os.Stat(outputPath); statErr != nil {
		errMsg := stderrText
		if errMsg == "" && execErr != nil {
			errMsg = execErr.Error()
		}
		if errMsg == "" {
			errMsg = "[BenchWire] No output produced, no error captured"
		}
		if writeErr := WriteFailureYaml(outputPath, errMsg, label); writeErr != nil {
			return runNumber, writeErr
		}
	}

	timerEnd := GetTime()
	cooldown, err := msleep(cooldownTimer)
	if err != nil {
		return runNumber, err
	}

	runNumber++
	progress.Tick(runNumber, timerStart, timerEnd, cooldown, specMatrixName)

	return runNumber, nil
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
		outputPreset, ok := outputPresets[specMatrix.Benchmarker]

		if !ok {
			return fmt.Errorf("[BenchWire] no output preset for benchmarker %q", specMatrix.Benchmarker)
		}
		
		// TODO: Turn this into a seperate function with preset support.
		if specMatrix.Benchmarker == "pyperf" {
			outputDir, _ = filepath.Split(outputDir)
			outputDir = filepath.Join(outputDir, "json")
		}

		for _, runTarget := range specMatrix.RunTargets{
			switch runTarget.Methodology {
				case "single":
					var err error
					for i := 1; i <= runTarget.Runs; i++ {
						runNumber, err = SingularRun(&progress, runNumber, outputPreset, outputDir, specMatrix.SpecMatrixName, runTarget.BinPath[0], runTarget.Flags[0], runTarget.Label[0], "", i, runTarget.CooldownTimer)
						if err != nil {
							return err
						}
					}

				case "sequential":
					var err error
					for i := 1; i <= runTarget.Runs; i++ {
						runNumber, err = SingularRun(&progress, runNumber, outputPreset, outputDir, specMatrix.SpecMatrixName, runTarget.BinPath[0], runTarget.Flags[0], runTarget.Label[0], "A", i, runTarget.CooldownTimer)
						if err != nil {
							return err
						}
					}

					for i := 1; i <= runTarget.Runs; i++ {
						runNumber, err = SingularRun(&progress, runNumber, outputPreset, outputDir, specMatrix.SpecMatrixName, runTarget.BinPath[1], runTarget.Flags[1], runTarget.Label[1], "B", i, runTarget.CooldownTimer)
						if err != nil {
							return err
						}
					}

				case "cycling":
					var err error
					for i := 1; i <= runTarget.Runs; i++ {
						runNumber, err = SingularRun(&progress, runNumber, outputPreset, outputDir, specMatrix.SpecMatrixName, runTarget.BinPath[0], runTarget.Flags[0], runTarget.Label[0], "A", i, runTarget.CooldownTimer)
						if err != nil {
							return err
						}

						runNumber, err = SingularRun(&progress, runNumber, outputPreset, outputDir, specMatrix.SpecMatrixName, runTarget.BinPath[1], runTarget.Flags[1], runTarget.Label[1], "B", i, runTarget.CooldownTimer)
						if err != nil {
							return err
						}
					}

				case "random interleaving":
					order := make([]byte, 0, runTarget.Runs*2)
					for i := 0; i < runTarget.Runs; i++ {
						order = append(order, 'A', 'B')
					}
					rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

					var err error
					countA, countB := 1, 1

					for _, side := range order {
						if side == 'A' {
							runNumber, err = SingularRun(&progress, runNumber, outputPreset, outputDir, specMatrix.SpecMatrixName, runTarget.BinPath[0], runTarget.Flags[0], runTarget.Label[0], "A", countA, runTarget.CooldownTimer)
							countA++
						} else {
							runNumber, err = SingularRun(&progress, runNumber, outputPreset, outputDir, specMatrix.SpecMatrixName, runTarget.BinPath[1], runTarget.Flags[1], runTarget.Label[1], "B", countB, runTarget.CooldownTimer)
							countB++
						}

						if err != nil {
							return err
						}
					}
				}
			}
	}
	fmt.Printf("\n[BenchWire] total: %.2fs\n", float64(GetTotalTimeSpent(progress))/1e9)
	return nil
}
