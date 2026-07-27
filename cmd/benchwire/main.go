package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"os/exec"
	"benchwire/internal/config"
	"benchwire/internal/runner"
)

func main() {
	scriptDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't get working directory:", err)
		os.Exit(1)
	}

	outputDir := filepath.Join(scriptDir, "results", "yaml")
	plotsDir := filepath.Join(scriptDir, "results", "plots")

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "couldn't create output dir:", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(plotsDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "couldn't create plots dir:", err)
		os.Exit(1)
	}

	fmt.Println("Warning: this will delete your previous runs' yaml files")
	fmt.Print("1 to use config.yaml | 2 for custom command: ")

	reader := bufio.NewReader(os.Stdin)
	option, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't read input:", err)
		os.Exit(1)
	}
	option = strings.TrimSpace(option)

	switch option {
		case "1":
			fmt.Println("[Config mode enabled]")
			cfg, err := config.LoadYamlConfig(scriptDir)
			if err != nil{
				fmt.Fprintln(os.Stderr, "couldn't load config:", err)
				os.Exit(1)
			}

			opts := runner.RunArgs {
				Methodology: cfg.Methodology,
				Runs:        cfg.Runs,
				CooldownMs:  cfg.CooldownTimer,
				OutputDir:   outputDir,
				Targets:     cfg.Targets,
			}
			
			if err := runner.Run(opts); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			pyCmd := exec.Command("python3", filepath.Join(scriptDir, "analyze.py"))
			
			pyCmd.Stdout = os.Stdout
			pyCmd.Stderr = os.Stderr
			if err := pyCmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

		case "2":
			fmt.Println("[custom mode - currently disabled]")
		default:
			fmt.Fprintln(os.Stderr, "Invalid option, exiting.")
			os.Exit(1)
	}
}
