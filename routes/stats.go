package routes

import (
	"net/http"
	"strings"
	"sync"
	_ "time/tzdata"

	"github.com/gin-gonic/gin"
	"github.com/oddmario/systemstats-agent/stats"
	"github.com/oddmario/systemstats-agent/utils"
	"github.com/tidwall/sjson"
)

func Stats(c *gin.Context) {
	if len(utils.AuthKey) > 0 {
		apiKey := strings.TrimSpace(c.Query("apikey"))

		if apiKey != utils.AuthKey {
			c.Status(http.StatusUnauthorized)

			return
		}
	}

	var wg sync.WaitGroup

	memoryStats, _ := stats.GetMemoryUsage()
	hostInfo, _ := stats.GetHostInfo()
	var networkStats []map[string]interface{} = nil
	var diskStats []map[string]interface{} = nil
	var averageTotalCpuUsagePercent float64 = 0.0
	var cpuUsagePercentPerCore []float64 = []float64{}
	var publicIPv4 string = ""
	var publicIPv6 string = ""

	wg.Add(1)
	go func() {
		defer wg.Done()

		networkStats, _ = stats.GetNetworks()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		diskStats, _ = stats.GetDisks()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		averageTotalCpuUsagePercent, _ = stats.GetCPUUsage()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		cpuUsagePercentPerCore, _ = stats.GetCPUUsagePerCore()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		publicIPv4, _ = stats.GetPublicIP(false)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		publicIPv6, _ = stats.GetPublicIP(true)
	}()

	wg.Wait()

	c.Header("Content-Type", "application/json")

	var i int64

	resJson, _ := sjson.Set("", "error", false)
	resJson, _ = sjson.Set(resJson, "host", map[string]interface{}{})
	resJson, _ = sjson.Set(resJson, "network", []map[string]interface{}{})
	resJson, _ = sjson.Set(resJson, "disk", []map[string]interface{}{})
	resJson, _ = sjson.Set(resJson, "memory", map[string]interface{}{})
	resJson, _ = sjson.Set(resJson, "public_ipv4", publicIPv4)
	resJson, _ = sjson.Set(resJson, "public_ipv6", publicIPv6)
	resJson, _ = sjson.Set(resJson, "average_cpu_usage_percent", averageTotalCpuUsagePercent)
	resJson, _ = sjson.Set(resJson, "cpu_usage_percent_per_core", cpuUsagePercentPerCore)

	i = 0
	for _, interface_ := range networkStats {
		var json_key_prefix string = "network." + utils.I64ToStr(i) + "."

		resJson, _ = sjson.Set(resJson, json_key_prefix+"interface_name", interface_["interface_name"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"mtu", interface_["mtu"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"flags", interface_["flags"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"hardware_address", interface_["hardware_addr"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"addresses", interface_["addrs"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"downrate_KB_per_sec", interface_["down_KB_per_sec"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"uprate_KB_per_sec", interface_["up_KB_per_sec"])

		i++
	}

	i = 0
	for _, disk := range diskStats {
		var json_key_prefix string = "disk." + utils.I64ToStr(i) + "."

		resJson, _ = sjson.Set(resJson, json_key_prefix+"device_name", disk["device"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"mountpoint", disk["mountpoint"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"mountopts", disk["mountopts"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"fstype", disk["fstype"])

		resJson, _ = sjson.Set(resJson, json_key_prefix+"free_bytes", disk["free_bytes"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"used_bytes", disk["used_bytes"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"total_size", disk["total_size"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"free_inodes", disk["free_inodes"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"used_inodes", disk["used_inodes"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"total_inodes", disk["total_inodes"])

		resJson, _ = sjson.Set(resJson, json_key_prefix+"read_speed_KB_per_sec", disk["read_speed"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"write_speed_KB_per_sec", disk["write_speed"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"iops", disk["iops"])

		i++
	}

	resJson, _ = sjson.Set(resJson, "memory.total_mem", memoryStats["total_mem"])
	resJson, _ = sjson.Set(resJson, "memory.used_mem", memoryStats["used_mem"])
	resJson, _ = sjson.Set(resJson, "memory.total_swap", memoryStats["total_swap"])
	resJson, _ = sjson.Set(resJson, "memory.used_swap", memoryStats["used_swap"])
	resJson, _ = sjson.Set(resJson, "memory.cached_memory", memoryStats["cached_memory"])

	resJson, _ = sjson.Set(resJson, "host.hostname", hostInfo["hostname"])
	resJson, _ = sjson.Set(resJson, "host.uptime_secs", hostInfo["uptime_secs"])
	resJson, _ = sjson.Set(resJson, "host.boot_timestamp", hostInfo["boot_timestamp"])
	resJson, _ = sjson.Set(resJson, "host.total_procs", hostInfo["total_procs"])
	resJson, _ = sjson.Set(resJson, "host.os", hostInfo["os"])
	resJson, _ = sjson.Set(resJson, "host.platform", hostInfo["platform"])
	resJson, _ = sjson.Set(resJson, "host.platform_family", hostInfo["platform_family"])
	resJson, _ = sjson.Set(resJson, "host.platform_version", hostInfo["platform_version"])
	resJson, _ = sjson.Set(resJson, "host.kernel_version", hostInfo["kernel_version"])
	resJson, _ = sjson.Set(resJson, "host.kernel_arch", hostInfo["kernel_arch"])
	resJson, _ = sjson.Set(resJson, "host.vz_system", hostInfo["vz_system"])
	resJson, _ = sjson.Set(resJson, "host.vz_role", hostInfo["vz_role"])
	resJson, _ = sjson.Set(resJson, "host.host_uuid", hostInfo["host_uuid"])
	resJson, _ = sjson.Set(resJson, "host.avg_load_1", hostInfo["avg_load_1"])
	resJson, _ = sjson.Set(resJson, "host.avg_load_5", hostInfo["avg_load_5"])
	resJson, _ = sjson.Set(resJson, "host.avg_load_15", hostInfo["avg_load_15"])
	resJson, _ = sjson.Set(resJson, "host.temperatures", []map[string]interface{}{})

	i = 0
	for _, temp_ := range hostInfo["temperatures"].([]map[string]interface{}) {
		var json_key_prefix string = "host.temperatures." + utils.I64ToStr(i) + "."

		resJson, _ = sjson.Set(resJson, json_key_prefix+"sensor_key", temp_["sensor_key"])
		resJson, _ = sjson.Set(resJson, json_key_prefix+"temperature", temp_["temperature"])

		i++
	}

	c.Writer.WriteString(utils.PrettifyJson(resJson))
}
