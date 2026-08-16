package resolver

import (
	"benchwire/internal/config"
	"log"
)

type ResolvedTarget struct {
	Label string
	BinPath string
	Flags []string 
}

type SpecMatrixTargets struct {
	Name string
	Benchmarker string 
	Methodology string 
	Runs int
	CooldownTimer config.CooldownTimer
	ResolvedTarget []ResolvedTarget
	Sequence config.Sequence
}

type FullPreset struct {
	Name string
	Flags []string
}

func Resolve(cfg config.YamlConfig) []SpecMatrixTargets {
	mode := cfg.Context.Mode
	if mode == "" {
		mode = "default"
	}

	defaultMethodology := cfg.Default.Methodology
	defaultRuns := cfg.Default.Runs
	defaultCooldownTimer := cfg.Default.CooldownTimer

	parentPresetName := make(map[string][]string, len(cfg.Presets))

	for _, preset := range cfg.Presets {
		parentPresetName[preset.Name] = preset.Flags
	}

	specMatrixTargets := make([]SpecMatrixTargets, len(cfg.SpecMatrix))
	
	// specMatrix default inheritence builder. 
	for i, spec := range cfg.SpecMatrix {
		specMatrixTargets[i].Name = spec.Name
		specMatrixTargets[i].Benchmarker = spec.Benchmarker

		// Populate with default settings if not overridden.
		if len(spec.Methodology) == 0 {
			spec.Methodology = defaultMethodology
		}
		specMatrixTargets[i].Methodology = spec.Methodology

		if spec.Runs == 0 {
			spec.Runs = defaultRuns
		}
		specMatrixTargets[i].Runs = spec.Runs

		if len(spec.CooldownTimer.Value) == 0 {
			spec.CooldownTimer = defaultCooldownTimer
		}
		specMatrixTargets[i].CooldownTimer = spec.CooldownTimer
		
		combinedFlags := append([]string{}, cfg.GlobalFlags...)
		combinedFlags = append(combinedFlags, spec.LocalFlags...)
		
		fullPreset := make([]FullPreset, len(spec.Presets))

		localPresetName := make(map[string][]string, len(cfg.Presets))
		for name, flags := range parentPresetName {
			localPresetName[name] = flags
		}
		
		// Create presets for this specMatrix.
		for j, specPresets := range spec.Presets {
			
			fullPreset[j].Name = specPresets.Name

			if len(specPresets.Inherit) > 0 {
				for _, presetName := range specPresets.Inherit {
					if inheritedFlags, ok := parentPresetName[presetName]; ok {
						fullPreset[j].Flags = append(fullPreset[j].Flags, inheritedFlags...)
					} else {
						log.Fatalf("Error: Inherited preset %q not found", presetName)
					}
				}
			}

			fullPreset[j].Flags = append(fullPreset[j].Flags, specPresets.Flags...)
			localPresetName[fullPreset[j].Name] = fullPreset[j].Flags
		}
		
		// Populate target flags.
		resolvedTarget := make([]ResolvedTarget, len(spec.Targets))
		specMatrixTargets[i].ResolvedTarget = make([]ResolvedTarget, len(spec.Targets))

		for j, target := range spec.Targets{
			resolvedTarget[j].Label = target.Label
			resolvedTarget[j].BinPath = target.BinPath
			resolvedTarget[j].Flags = append(resolvedTarget[j].Flags, combinedFlags...)

			for _, preset := range target.Preset {
				resolvedTarget[j].Flags = append(resolvedTarget[j].Flags, localPresetName[preset]...)
			}

			resolvedTarget[j].Flags = append(resolvedTarget[j].Flags, target.Flags...)
			specMatrixTargets[i].ResolvedTarget[j] = resolvedTarget[j]
		}

		specMatrixTargets[i].Sequence = spec.Sequence
	}
	
	return specMatrixTargets 
} 
