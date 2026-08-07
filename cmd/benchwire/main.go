package main

import (
	"fmt"
	"os"
	"path/filepath"
	"os/exec"
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

	cfg, err := config.LoadYamlConfig(scriptDir)
	if err != nil{
		fmt.Fprintln(os.Stderr, "couldn't load config:", err)
		os.Exit(1)
	}
	
	// Once support for other benchmarkersr is added, updated allowedRoot.
	allowedRoot := outputDir
	if len(os.Args) > 1 && os.Args[1] == "clean" {
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

	pyCmd := exec.Command("python3", filepath.Join(scriptDir, "analyze.py"))
	
	pyCmd.Stdout = os.Stdout
	pyCmd.Stderr = os.Stderr
	if err := pyCmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
