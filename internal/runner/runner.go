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

func buildArgs(flags []string, preset OutputPreset, outputDir, matrixName,
			   label, prefix string, run int) ([]string, string) {
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

type RunContext struct {
	OutputPreset   OutputPreset
	OutputDir      string
	SpecMatrixName string
	BinPath        string
	Flags          []string
	Label          string
	Prefix         string
	CooldownTimer  config.CooldownTimer
	TotalRuns      int
}

func checkError(outputPath string) bool {
	data, err := os.ReadFile(outputPath)
	if err != nil {
		return false
	}

	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "error:") {
			continue
		}
		rest := strings.TrimSpace(line[len("error:"):])
		return rest != "" && rest != "''"
	}

	return false
}

func SingularRun(progress *ProgressBar, runNumber int, ctx RunContext, run int) (int, bool, error) {
	timerStart := GetTime()
	args, outputPath := buildArgs(ctx.Flags, ctx.OutputPreset, ctx.OutputDir, ctx.SpecMatrixName, ctx.Label, ctx.Prefix, run)
	os.Remove(outputPath)
	stderrText, execErr := execute(ctx.OutputPreset.Execution, ctx.BinPath, args)

	if _, statErr := os.Stat(outputPath); statErr != nil {
		errMsg := stderrText
		if errMsg == "" && execErr != nil {
			errMsg = execErr.Error()
		}

		if errMsg == "" {
			errMsg = "[BenchWire] No output produced, no error captured"
		}

		if writeErr := WriteFailureYaml(outputPath, errMsg, ctx.Label); writeErr != nil {
			return runNumber, true, writeErr
		}

		if run == 1 && ctx.TotalRuns > 10 {
			return runNumber, true, nil
		}
	} else if run == 1 && ctx.TotalRuns > 10 {
		if checkError(outputPath) {
			return runNumber, true, nil
		}
	}

	timerEnd := GetTime()
	cooldown, err := msleep(ctx.CooldownTimer)
	if err != nil {
		return runNumber, false, err
	}

	runNumber++
	progress.Tick(runNumber, timerStart, timerEnd, cooldown, ctx.SpecMatrixName)

	return runNumber, false, nil
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
		
		var skip bool = false
		for _, runTarget := range specMatrix.RunTargets{
			var err error

			if runTarget.Methodology == "single" {
				ctx := RunContext{
					OutputPreset: outputPreset, OutputDir: outputDir,
					SpecMatrixName: specMatrix.SpecMatrixName, 
					BinPath: runTarget.BinPath[0], Flags: runTarget.Flags[0],
					Label: runTarget.Label[0], Prefix: "",
					CooldownTimer: runTarget.CooldownTimer, TotalRuns: runTarget.Runs,
				}

				for i := 1; i <= runTarget.Runs; i++ {
					runNumber, skip, err = SingularRun(&progress, runNumber, ctx, i)
					if err != nil {
						return err
					}
					
					if skip {
						runNumber += runTarget.Runs
						skip = false
						break
					}
				}
				continue
			}

			ctxA := RunContext{
				OutputPreset: outputPreset, OutputDir: outputDir,
				SpecMatrixName: specMatrix.SpecMatrixName,
				BinPath: runTarget.BinPath[0], Flags: runTarget.Flags[0],
				Label: runTarget.Label[0], Prefix: "A",
				CooldownTimer: runTarget.CooldownTimer, TotalRuns: runTarget.Runs,
			}

			ctxB := RunContext{
				OutputPreset: outputPreset, OutputDir: outputDir,
				SpecMatrixName: specMatrix.SpecMatrixName,
				BinPath: runTarget.BinPath[1], Flags: runTarget.Flags[1],
				Label: runTarget.Label[1], Prefix: "B",
				CooldownTimer: runTarget.CooldownTimer, TotalRuns: runTarget.Runs,
			}

			switch runTarget.Methodology {
				case "sequential":
					for i := 1; i <= runTarget.Runs; i++ {
						runNumber, skip, err = SingularRun(&progress, runNumber, ctxA, i)
						if err != nil {
							return err
						}

						if skip {
							runNumber += runTarget.Runs
							skip = false
							break
						}
					}

					for i := 1; i <= runTarget.Runs; i++ {
						runNumber, skip, err = SingularRun(&progress, runNumber, ctxB, i)
						if err != nil {
							return err
						}

						if skip {
							runNumber += runTarget.Runs
							skip = false
							break
						}
					}

				case "cycling":
					for i := 1; i <= runTarget.Runs; i++ {
						runNumber, skip, err = SingularRun(&progress, runNumber, ctxA, i)
						if err != nil {
							return err
						}

						runNumber, skip, err = SingularRun(&progress, runNumber, ctxB, i)
						if err != nil {
							return err
						}

						if skip {
							runNumber += runTarget.Runs
							skip = false
							break
						}
					}

				case "random interleaving":
					order := make([]byte, 0, runTarget.Runs*2)
					for i := 0; i < runTarget.Runs; i++ {
						order = append(order, 'A', 'B')
					}
					rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })

					countA, countB := 1, 1
					for _, side := range order {
						if side == 'A' {
							runNumber, skip, err = SingularRun(&progress, runNumber, ctxA, countA)
							countA++
						} else {
							runNumber, skip, err = SingularRun(&progress, runNumber, ctxB, countB)
							countB++
						}

						if err != nil {
							return err
						}

						if skip {
							runNumber+= runTarget.Runs
							skip = false
							break
						}
					}
				}
			}	
	}
	fmt.Printf("\n[BenchWire] total: %.2fs\n", float64(GetTotalTimeSpent(progress))/1e9)
	return nil
}
