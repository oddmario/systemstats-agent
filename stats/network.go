package stats

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/shirou/gopsutil/net"
)

func GetNetworks() ([]map[string]interface{}, error) {
	// Get initial stats
	stats1, err := net.IOCounters(true) // true = per interface stats
	if err != nil {
		return nil, err
	}

	// Wait for a second
	time.Sleep(time.Second)

	// Get stats after interval
	stats2, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	result := []map[string]interface{}{}

	// Calculate bandwidth for each interface
	for i, s1 := range stats1 {
		s2 := stats2[i]
		bytesSent := s2.BytesSent - s1.BytesSent
		bytesRecv := s2.BytesRecv - s1.BytesRecv

		var interface_ net.InterfaceStat = net.InterfaceStat{}
		var wasInterfaceFound bool = false

		for _, if_ := range interfaces {
			if if_.Name == s1.Name {
				interface_ = if_
				wasInterfaceFound = true

				break
			}
		}

		if !wasInterfaceFound {
			continue
		}

		var interfaceAddrs []string = []string{}
		for _, addr := range interface_.Addrs {
			interfaceAddrs = append(interfaceAddrs, addr.Addr)
		}

		result = append(result, map[string]interface{}{
			"interface_name":  s1.Name,
			"mtu":             interface_.MTU,
			"flags":           interface_.Flags,
			"hardware_addr":   interface_.HardwareAddr,
			"addrs":           interfaceAddrs,
			"down_KB_per_sec": float64(bytesRecv) / 1024, // KB/s
			"up_KB_per_sec":   float64(bytesSent) / 1024, // KB/s
		})
	}

	return result, nil
}

func GetPublicIP(return_ipv6 bool) (string, error) {
	dnsServer := "ns1.google.com:53"
	queryDomain := "o-o.myaddr.l.google.com"
	if return_ipv6 {
		dnsServer = "[2001:4860:4802:32::a]:53"
	}

	msg := new(dns.Msg)
	msg.SetQuestion(queryDomain+".", dns.TypeTXT)

	res, err := dns.Exchange(msg, dnsServer)
	if err != nil {
		return "", err
	}

	if len(res.Answer) < 1 {
		return "", fmt.Errorf("no answer in dns response")
	}

	ip, ok := res.Answer[0].(*dns.TXT)
	if !ok {
		return "", fmt.Errorf("failed to parse dns response")
	}

	return ip.Txt[0], nil
}
