package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"os/exec"
	"benchwire/internal/builder"
	"benchwire/internal/config"
	"benchwire/internal/runner"
	"benchwire/internal/resolver"
	"benchwire/internal/scheduler"
	"benchwire/internal/cleanup"
)

func main() {
	scriptDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't get working directory:", err)
		os.Exit(1)
	}

	outputDir := filepath.Join(scriptDir, "results", "yaml")

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "couldn't create output dir:", err)
		os.Exit(1)
	}
	
	configPath := scriptDir
	configFile := "config.yaml"
	clean := false
	
	// TODO: If this get's too long, move it to cli.go or similar.
	for _, arg := range os.Args[1:] {
		switch {
			case arg == "--clean":
				clean = true

			case arg == "--ci": // Not implamented yet.
				configPath = "configs"
				configFile = "ci.yaml"

			case strings.HasPrefix(arg, "--config="):
				full := strings.TrimPrefix(arg, "--config=")
				configPath = filepath.Dir(full)
				configFile = filepath.Base(full)

			case strings.HasPrefix(arg, "--buildConfig="):
				loc := strings.TrimPrefix(arg, "--buildConfig=")
				loc, err = builder.FindConfig(filepath.Join(scriptDir, "configs"), loc)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				configPath, configFile, err = builder.ConfigBuilder(loc, scriptDir)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				runConfig := bufio.NewScanner(os.Stdin)
				shouldRun := true

				for {
					fmt.Print("[BenchWire] Run config? Y/N ")
					var yesNo string
					if runConfig.Scan() {
						yesNo = runConfig.Text()
					}

					if err := runConfig.Err(); err != nil {
						fmt.Fprintln(os.Stderr, "[BenchWire] reading standard input: ", err)
					}

					switch yesNo {
						case "y", "Y":
							shouldRun = true

						case "n", "N":
							shouldRun = false

						default:
							fmt.Println("[BenchWire] Please enter Y or N.")
							continue
					}
					break
				}

				if !shouldRun {
					return
				}
		}
	}

	cfg, err := config.LoadYamlConfig(configPath, configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[BenchWire] Couldn't load config:", err)
		os.Exit(1)
	}

	// Once support for other benchmarkersr is added, updated allowedRoot.
	allowedRoot := outputDir
	if clean {
		if err := cleanup.SafeRemove(outputDir, allowedRoot); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("cleaned.")
		return
	}
	
	fmt.Println("[Benchwire] Activated")

	resolved := resolver.Resolve(*cfg)
	scheduled, err := scheduler.Schedule(resolved)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := runner.Run(scheduled, outputDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	pyCmd := exec.Command("python3", filepath.Join(scriptDir, "analyze.py"), filepath.Join(configPath, configFile))
	pyCmd.Stdout = os.Stdout
	pyCmd.Stderr = os.Stderr

	if err := pyCmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
