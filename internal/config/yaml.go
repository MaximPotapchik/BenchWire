package config

import (
	//"bufio"
	"fmt"
	"os"
	"path/filepath"
	"go.yaml.in/yaml/v3"
)

type YamlConfig struct {
	Methodology string `yaml:"methodology"`
	Runs int `yaml:"runs"`
	CooldownTimer int `yaml:"cooldownTimer"`
	Targets []Target `yaml:"targets"`
}

type Target struct {
	Label string `yaml:"label"`
	BinPath string `yaml:"binPath"`
	Flags []string `yaml:"flags"`
}

func LoadYamlConfig(targetDir string) (*YamlConfig, error) {
	path := filepath.Join(targetDir, "config.yaml")
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("no config.yaml found, copy exampleconfig.yaml to config.yaml and fill it in")
	}

	file, err := os.Open(path)
    if err != nil {
        return  nil, err
    }
    defer file.Close()

	var data YamlConfig

	err = yaml.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	for i := range data.Targets {
		data.Targets[i].BinPath = os.ExpandEnv(data.Targets[i].BinPath)
	}
	
	return &data, nil
}
