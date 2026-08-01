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
	CooldownTimer CooldownTimer `yaml:"cooldownTimer"`
	Targets []Target `yaml:"targets"`
}

type Target struct {
	Label string `yaml:"label"`
	BinPath string `yaml:"binPath"`
	Flags []string `yaml:"flags"`
}

type CooldownTimer struct {
	Value      string
	Randomized bool
}

func (c *CooldownTimer) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
		case yaml.ScalarNode:
			var raw string
			if err := value.Decode(&raw); err != nil {
				return err
			}

			c.Value = raw
			c.Randomized = false
			return nil

		case yaml.MappingNode:
			var obj struct {
				RandomizeWithin int    `yaml:"randomizeWithin"`
				Precision       string `yaml:"precision"`
			}

			if err := value.Decode(&obj); err != nil {
				return err
			}

			c.Value = fmt.Sprintf("%d%s", obj.RandomizeWithin, obj.Precision)
			c.Randomized = true
			return nil
	}
	return fmt.Errorf("cooldownTimer: unsupported yaml syntax.")
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
