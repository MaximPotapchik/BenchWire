package config

import (
	//"bufio"
	"fmt"
	"os"
	"path/filepath"
	"go.yaml.in/yaml/v3"
)

// conig structs are followed by their unmarshalling method.

type YamlConfig struct {
	Context Context `yaml:"context"`
	Default Default `yaml:"default"`
	GlobalFlags []string `yaml:"globalFlags"`
	Presets []Preset `yaml:"presets"`
	SpecMatrix []SpecMatrix `yaml:"specMatrix"`
	MatrixSequence []string `yaml:"matrixSequence"`
	Analysis Analysis `yaml:"analysis"`
}

type Context struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
	//Schedule Schedule `yaml:"schedule"`
}

type Default struct {
	Methodology string `yaml:"methodology"`
	Runs int `yaml:"runs"`
	CooldownTimer CooldownTimer `yaml:"cooldownTimer"`	
}

type SpecMatrix struct {
	Name string `yaml:"name"`
	Benchmarker string `yaml:"benchmarker"`
	Methodology string `yaml:"methodology"`
	Runs int `yaml:"runs"`
	CooldownTimer CooldownTimer `yaml:"cooldownTimer"`
	LocalFlags []string `yaml:"localFlags"`
	Presets []Preset `yaml:"presets"`
	Targets []Target `yaml:"targets"`
	Sequence Sequence `yaml:"sequence"`
}

type Preset struct {
	Name string `yaml:"name"`
	Inherit []string `yaml:"inherit,omitempty"`
	Flags []string `yaml:"flags"`
}

type Target struct {
	Label string `yaml:"label"`
	BinPath string `yaml:"binPath"`
	Preset []string `yaml:"preset,omitempty"`
	Flags []string `yaml:"flags"`
}

type CooldownTimer struct {
	Value string
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
				RandomizeWithin int `yaml:"randomizeWithin"`
				Precision string `yaml:"precision"`
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

func (c CooldownTimer) MarshalYAML() (interface{}, error) {
	return c.Value, nil
}

// Sequence parsing for whether it is single or comparison. 
type Sequence [][]string

func (s *Sequence) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return fmt.Errorf("sequence: expected a list")
	}
	steps := make([][]string, 0, len(value.Content))
	for _, item := range value.Content {
		switch item.Kind {
			case yaml.ScalarNode:
				var label string
				if err := item.Decode(&label); err != nil {
					return err
				}
				steps = append(steps, []string{label})
			case yaml.SequenceNode:
				var labels []string
				if err := item.Decode(&labels); err != nil {
					return err
				}
				steps = append(steps, labels)
			default:
				return fmt.Errorf("sequence step: unsupported yaml syntax")
		}
	}
	*s = steps
	return nil
}

type Analysis struct {
	Output AnalysisOutput `yaml:"output"`
}

type AnalysisOutput struct {
	DirName string `yaml:"dirName"`
	StoreIn string `yaml:"storeIn"`
	Per string `yaml:"per"`
	ShowErrored bool `yaml:"showErrored"`
	Plots Plots `yaml:"plots"`
}

type Plots struct {
	Disable bool `yaml:"disable"`
}

func LoadYamlConfig(targetDir string, filename ...string) (*YamlConfig, error) {
	name := "config.yaml"
	if len(filename) > 0 {
		name = filename[0]
	}
	path := filepath.Join(targetDir, name)
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("no config found for path: %s", path)
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
	
	for i := range data.SpecMatrix {
		for j := range data.SpecMatrix[i].Targets {
			data.SpecMatrix[i].Targets[j].BinPath = os.ExpandEnv(data.SpecMatrix[i].Targets[j].BinPath)
		}
	}

	return &data, nil
}
