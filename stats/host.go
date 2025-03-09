package stats

import (
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
)

func GetHostInfo() (map[string]interface{}, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}

	temp_stat, err := host.SensorsTemperatures()
	if err != nil {
		return nil, err
	}

	load, err := load.Avg()
	if err != nil {
		return nil, err
	}

	temp_stat_final := []map[string]interface{}{}

	for _, temp_stat_item := range temp_stat {
		temp_stat_final = append(temp_stat_final, map[string]interface{}{
			"sensor_key":  temp_stat_item.SensorKey,
			"temperature": temp_stat_item.Temperature,
		})
	}

	return map[string]interface{}{
		"hostname":         info.Hostname,
		"uptime_secs":      info.Uptime,
		"boot_timestamp":   info.BootTime,
		"total_procs":      info.Procs,
		"os":               info.OS,
		"platform":         info.Platform,
		"platform_family":  info.PlatformFamily,
		"platform_version": info.PlatformVersion,
		"kernel_version":   info.KernelVersion,
		"kernel_arch":      info.KernelArch,
		"vz_system":        info.VirtualizationSystem,
		"vz_role":          info.VirtualizationRole,
		"host_uuid":        info.HostID,
		"avg_load_1":       load.Load1,
		"avg_load_5":       load.Load5,
		"avg_load_15":      load.Load15,
		"temperatures":     temp_stat_final,
	}, nil
}
