package iptrie

import (
	"net"
	"sort"
	"testing"
)

func TestIpTrie_FindByIp(t *testing.T) {
	ipList := []string{
		"192.168.100.2",
		"192.168.10.10-192.168.20.9",
		"172.16.100.1/24",
		"2001:db8::68",
	}

	ipTrie := NewTrie()
	for _, ip := range ipList {
		ipTrie.InsertIpAddr(NewIpAddr(ip), nil)
	}

	ipResult := map[string]bool{
		"192.168.100.2":  true,
		"192.168.100.1":  false,
		"192.168.10.10":  true,
		"192.168.19.255": true,
		"192.168.20.10":  false,
		"172.16.100.100": true,
		"172.16.101.100": false,
		"2001:db8::68":   true,
		"2001:db8::70":   false,
	}

	for ip, result := range ipResult {
		if ipTrie.FindByIp(net.ParseIP(ip)) != result {
			t.Errorf("%s should be %v", ip, result)
		}
	}
}

func TestIpTrie_FindLowestByIp(t *testing.T) {
	ipValues := map[string]int{
		"192.168.100.0/24":              1,
		"192.168.100.1-192.168.100.100": 2,
	}

	var ipAddrList []*IpAddr
	for ip, _ := range ipValues {
		ipAddrList = append(ipAddrList, NewIpAddr(ip))
	}

	sort.Slice(ipAddrList, func(i, j int) bool {
		// sort asc
		return ipAddrList[i].Count().Cmp(ipAddrList[j].Count()) < 0
	})

	ipTrie := NewTrie()
	for _, ipAddr := range ipAddrList {
		ipTrie.InsertIpAddr(ipAddr, ipValues[ipAddr.String()])
	}

	ipResult := map[string]int{
		"192.168.100.1":   2,
		"192.168.100.101": 1,
	}

	for ip, result := range ipResult {
		if ipTrie.FindLowestByIp(net.ParseIP(ip)) != result {
			t.Errorf("%s should be %v", ip, result)
		}
	}
}
