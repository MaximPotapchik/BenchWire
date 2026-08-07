package scheduler

import (
	"benchwire/internal/resolver"
	"benchwire/internal/config"
	"fmt"
)

type ScheduledMatrix struct {
	SpecMatrixName string
	Benchmarker string
	RunTargets []RunTarget
}

type RunTarget struct {
	Number int
	Methodology string
	Runs int
	CooldownTimer config.CooldownTimer
	Label []string
	BinPath []string
	Flags [][]string
}

// Hash lookup helper.
func lookupTarget(targetHash map[string]resolver.ResolvedTarget, label string) (resolver.ResolvedTarget, error) {
	target, ok := targetHash[label]
	if !ok {
		return resolver.ResolvedTarget{}, fmt.Errorf("sequence references unknown label %q", label)
	}
	return target, nil
}

func Schedule(resolvedCfg []resolver.SpecMatrixTargets) ([]ScheduledMatrix, error){

	scheduledMatrix := make([]ScheduledMatrix, len(resolvedCfg))
	
	runNumber := 0

	for i, spec := range resolvedCfg {
		scheduledMatrix[i].SpecMatrixName = resolvedCfg[i].Name	
		scheduledMatrix[i].Benchmarker = resolvedCfg[i].Benchmarker

		targetHash := make(map[string]resolver.ResolvedTarget, len(spec.ResolvedTarget))

		for _, target := range spec.ResolvedTarget {
			targetHash[target.Label] = target
		}
		
		for _, order := range spec.Sequence {
			runNumber++
			runTarget := RunTarget{
				Number: runNumber,
				Methodology: spec.Methodology,
				Runs: spec.Runs,
				CooldownTimer: spec.CooldownTimer,
			}

			if len(order) > 1 {

				a, err := lookupTarget(targetHash, order[0])
				if err != nil {
					return nil, err
				}
				b, err := lookupTarget(targetHash, order[1])
				if err != nil {
					return nil, err
				}

				runTarget.Label = []string{a.Label, b.Label}
				runTarget.BinPath = []string{a.BinPath, b.BinPath}
				runTarget.Flags = [][]string{a.Flags, b.Flags}

			} else if len(order) == 1 {

				a, err := lookupTarget(targetHash, order[0])
				if err != nil {
					return nil, err
				}

				runTarget.Methodology = "single"
				runTarget.Label = []string{a.Label}
				runTarget.BinPath = []string{a.BinPath}
				runTarget.Flags = [][]string{a.Flags}

			} else {
				return nil, fmt.Errorf("[Benchwire] Error: Invalid length %d", len(order))
			}

			scheduledMatrix[i].RunTargets = append(scheduledMatrix[i].RunTargets, runTarget)
		}
	}
	return scheduledMatrix, nil
}
