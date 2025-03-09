package stats

import "github.com/shirou/gopsutil/mem"

func GetMemoryUsage() (map[string]uint64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return map[string]uint64{
		"total_mem":     v.Total,
		"used_mem":      v.Used,
		"total_swap":    v.SwapTotal,
		"used_swap":     (v.SwapTotal - v.SwapFree),
		"cached_memory": v.Cached,
	}, nil
}
