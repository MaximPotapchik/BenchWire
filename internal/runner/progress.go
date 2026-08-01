package runner

import (
	"time"
	"strings"
	"fmt"
)

type ProgressBar struct {
	TotalRuns int
	RunNumber int
	CapturedTimes []RunTiming
	RunningSum int64
	CooldownSum int64
	EstimatedETA int
	Percent int
}

type RunTiming struct {
	Run int
	TimerStart int64
	TimerEnd int64
	Reading int64
}

func (p *ProgressBar) Tick(run int, start, end int64, cooldownNs int64) {
	reading := UpdateProgress(run, start, end)
	p.RunNumber = run
	p.RunningSum += reading
	p.CapturedTimes = append(p.CapturedTimes, RunTiming{Run: run, TimerStart: start, TimerEnd: end, Reading: reading})

	// TODO: This can be improved. Auto-scale to each any unit prefix.
	runningAvg := GetRunningAverage(p.RunningSum, int64(run))
	p.Percent = (100 * run) / p.TotalRuns

	// Progress bar 
	filledRatio := p.Percent / 5
	filled := strings.Repeat("=", filledRatio)
	empty := strings.Repeat(" ", 20-filledRatio)
	
	// ETA calculations
	p.CooldownSum += cooldownNs
    avgCooldownMs := float64(p.CooldownSum) / float64(run) / 1000000
    p.EstimatedETA = (p.TotalRuns - run) * int(runningAvg + avgCooldownMs) / 1000

	fmt.Printf("[%s%s] %d%% (%d/%d) | avg %.2fms/run | ETA %ds \r", filled, empty, p.Percent, run, p.TotalRuns, runningAvg, p.EstimatedETA)
}

func GetTime() int64 {
	return time.Now().UnixNano()
}

func CaptureRunTime(start int64, end int64) int64 {
	return end - start
}

func GetRunningAverage(runningSum int64, runNumber int64) float64 {
	return float64(runningSum) / float64(runNumber) / 1000000
}

// TODO: Maybe return a struct instead so that CapturedTimes can store.
func UpdateProgress(run int, startTime int64, endTime int64) int64 {
	timing := RunTiming {
		Run: run,
		TimerStart: startTime,
		TimerEnd: endTime,
		Reading: 0,
	}

	timing.Reading = CaptureRunTime(timing.TimerStart, timing.TimerEnd)

	return timing.Reading
}

func GetTotalTimeSpent(bar ProgressBar) int64 {
	return bar.RunningSum + bar.CooldownSum
}
