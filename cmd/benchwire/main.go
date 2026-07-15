package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"os/exec"
	"benchwire/internal/config"
	"benchwire/internal/exegesis"
	"benchwire/internal/runners"
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
	fmt.Print("1 to use ENV | 2 for custom command: ")

	reader := bufio.NewReader(os.Stdin)
	option, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't read input:", err)
		os.Exit(1)
	}
	option = strings.TrimSpace(option)

	switch option {
		case "1":
			fmt.Println("[env mode enabled]")
			env, err := config.LoadEnv(scriptDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			methodology := os.Getenv("METHODOLOGY")
			formatted, err := exegesis.Format(methodology, env.RawFlags)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			runsN, _ := strconv.Atoi(os.Getenv("RUNS"))
			cooldownMs := 500
			if v, ok := os.LookupEnv("COOLDOWNTIMER"); ok { 
				if parsed, err := strconv.Atoi(v); err == nil {
					cooldownMs = parsed
				}
			}

			if err := runner.Run(methodology, runsN, cooldownMs, outputDir, formatted); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			labels := os.Getenv("LABEL")
			if labels == "" {
				labels = "Run"
			}

			if methodology != "single" {
				la, lb := os.Getenv("LABELA"), os.Getenv("LABELB")
				if la != "" && lb != "" {
					labels = la + "|" + lb
				} else {
					labels = "A|B"
				}
			}

			pyCmd := exec.Command("python3", filepath.Join(scriptDir, "analyze.py"),
				methodology,
				strconv.Itoa(runsN),
				strconv.Itoa(cooldownMs),
				labels,
			)
			
			pyCmd.Stdout = os.Stdout
			pyCmd.Stderr = os.Stderr
			if err := pyCmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

		case "2":
			fmt.Println("[custom mode]")
		default:
			fmt.Fprintln(os.Stderr, "Invalid option, exiting.")
			os.Exit(1)
	}
}
