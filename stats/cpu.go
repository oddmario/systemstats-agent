package stats

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// GetCPUUsage returns average CPU usage across all cores in percentage
func GetCPUUsage() (float64, error) {
	// Get per-core percentages
	percentages, err := cpu.Percent(time.Second, true)
	if err != nil {
		return 0, err
	}

	// Calculate average across all cores
	var total float64
	for _, percent := range percentages {
		total += percent
	}
	avg := total / float64(len(percentages))
	return avg, nil
}

// GetCPUUsagePerCode returns average CPU usage per core in percentages
func GetCPUUsagePerCore() ([]float64, error) {
	return cpu.Percent(time.Second, true)
}

// GetDetailedCPUUsage provides more comprehensive CPU metrics in percentages
func GetDetailedCPUUsage() (map[string]float64, error) {
	// Get CPU times for more detailed breakdown
	times, err := cpu.Times(false) // false = total, not per-CPU
	if err != nil {
		return nil, err
	}

	// Take initial measurement
	initial := times[0]

	// Wait for sample duration
	time.Sleep(time.Second)

	// Take final measurement
	times, err = cpu.Times(false)
	if err != nil {
		return nil, err
	}
	final := times[0]

	// Calculate deltas
	totalTime := final.Total() - initial.Total()
	if totalTime == 0 {
		return nil, fmt.Errorf("no cpu time elapsed")
	}

	results := make(map[string]float64)
	results["user"] = ((final.User - initial.User) / totalTime) * 100
	results["system"] = ((final.System - initial.System) / totalTime) * 100
	results["idle"] = ((final.Idle - initial.Idle) / totalTime) * 100
	results["total"] = 100 - results["idle"]

	return results, nil
}
