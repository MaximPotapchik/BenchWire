package builder

import (
	"benchwire/internal/config"
	"errors"
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"go.yaml.in/yaml/v3"
)

func BuildPaths(cfg config.YamlConfig, binpath string) *config.YamlConfig {
	pathOne, pathTwo, _ := strings.Cut(binpath, ",")

	for i, specMatrix := range cfg.SpecMatrix {
		for j, target := range specMatrix.Targets {
			if target.BinPath == "Bin1" {
				cfg.SpecMatrix[i].Targets[j].BinPath = pathOne
			} else if target.BinPath == "Bin2" {
				cfg.SpecMatrix[i].Targets[j].BinPath = pathTwo
			}
		}
	}

	return &cfg
}

func BuildConfig(cfg config.YamlConfig, mode string, binpath string) *config.YamlConfig {
	switch mode {
		case "default":
			BuildPaths(cfg, binpath)
	}

	return &cfg
}

func FindConfig(configsDir, userInput string) (string, error) {
	target := filepath.Base(userInput)
	var matches []string

	err := filepath.WalkDir(configsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && d.Name() == target {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("[BenchWire] no config named %q found anywhere under %s", target, configsDir)
	}

	if len(matches) == 1 {
		return matches[0], nil
	}

	return "", fmt.Errorf("[BenchWire] %q found in multiple places:  %s", target, strings.Join(matches, "\n "))
	
}

// TODO: - Needs a system to allow custom data to be passed. Via yaml, json, csv, etc.
// - Add a mode that combines two different config's specMatrixes.
func ConfigBuilder(configLoc string, scriptDir string) (string, string, error) {
	
	fmt.Print("[BenchWire] Initiated config builder. For more information, please look at docs/commands.md \n")

	dir := filepath.Dir(configLoc)
	file := filepath.Base(configLoc)

	if file == "sweep" {
		dir, _ = BuildSweep(scriptDir, "latency")
		file += ".yaml"
		return dir, file, nil
	}

	cfg, err := config.LoadYamlConfig(dir, file)
	if err != nil {
        fmt.Fprintf(os.Stderr, "[BenchWire] failed to load YAML config: %v\n", err)
		os.Exit(1)
    }

	var newConfigName string
	fmt.Print("[BenchWire] New config name: ")
	fmt.Scanln(&newConfigName)

	cfg.Context.Name = newConfigName
	
	var binpath string
	fmt.Print("[BenchWire] Binary path(s). Comma seperated if > 1: ")
	fmt.Scanln(&binpath)
	
	BuildConfig(*cfg, "default", binpath)
	builtConfig, _ := yaml.Marshal(&cfg)
	
	custConfigDir := filepath.Join(scriptDir, "configs/custom")
	err = os.Mkdir(custConfigDir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		fmt.Fprintln(os.Stderr, "[BenchWire] couldn't create custom config directory. ", err)
		os.Exit(1)
	}

	custConfigDir = filepath.Join(custConfigDir, newConfigName)
	custConfigDir += ".yaml"
	err = os.WriteFile(custConfigDir, builtConfig, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[BenchWire] couldn't create custom config. ", err)
	}
	
	dir = filepath.Dir(custConfigDir)
	file = filepath.Base(custConfigDir)

	return dir, file, err
}
