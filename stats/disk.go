package stats

import (
	"strings"
	"time"

	"github.com/shirou/gopsutil/disk"
)

func GetDisks() ([]map[string]interface{}, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	deviceNames := []string{}

	for _, partition := range partitions {
		if !strings.HasPrefix(partition.Device, "/dev/") {
			continue
		}

		deviceNames = append(deviceNames, partition.Device)
	}

	if len(deviceNames) == 0 {
		return nil, nil
	}

	// Get initial I/O stats
	stats1, err := disk.IOCounters(deviceNames...)
	if err != nil {
		return nil, err
	}

	// Wait for a second
	time.Sleep(time.Second)

	// Get final I/O stats
	stats2, err := disk.IOCounters(deviceNames...)
	if err != nil {
		return nil, err
	}

	result := []map[string]interface{}{}

	for _, partition := range partitions {
		if !strings.HasPrefix(partition.Device, "/dev/") {
			continue
		}

		if len(partition.Mountpoint) <= 0 {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		// Calculate I/O speeds
		deviceKey := strings.TrimPrefix(partition.Device, "/dev/")
		s1, ok1 := stats1[deviceKey]
		s2, ok2 := stats2[deviceKey]
		var readSpeed, writeSpeed, iops float64

		if ok1 && ok2 {
			readBytes := float64(s2.ReadBytes - s1.ReadBytes)
			writeBytes := float64(s2.WriteBytes - s1.WriteBytes)
			ioCount := float64((s2.ReadCount + s2.WriteCount) - (s1.ReadCount + s1.WriteCount))

			readSpeed = readBytes / 1024   // KB/s (1 second sample)
			writeSpeed = writeBytes / 1024 // KB/s (1 second sample)
			iops = ioCount                 // ops per second (1 second sample)
		}

		result = append(result, map[string]interface{}{
			"device":       partition.Device,
			"mountpoint":   partition.Mountpoint,
			"mountopts":    partition.Opts,
			"fstype":       partition.Fstype,
			"free_bytes":   usage.Free,
			"used_bytes":   usage.Used,
			"total_size":   usage.Total,
			"free_inodes":  usage.InodesFree,
			"used_inodes":  usage.InodesUsed,
			"total_inodes": usage.InodesTotal,
			"read_speed":   readSpeed,
			"write_speed":  writeSpeed,
			"iops":         iops,
		})
	}

	return result, nil
}
