package iptrie

import (
	"math/big"
	"net"
	"strings"
)

const (
	IpaddrKindCidr = iota
	IpaddrKindRange
)

type IpAddr struct {
	Kind  byte
	Ip    net.IP
	IpNet *net.IPNet
	Ip1   net.IP
	Ip2   net.IP
	count *big.Int
}

func NewIpAddr(s string) *IpAddr {
	if strings.Contains(s, "-") {
		ips := strings.Split(s, "-")
		if len(ips) == 2 {
			ip1, ip2 := net.ParseIP(ips[0]), net.ParseIP(ips[1])
			if ip1 != nil && ip2 != nil && compareIP(ip1, ip2) == -1 {
				return &IpAddr{
					Kind: IpaddrKindRange,
					Ip1:  ip1,
					Ip2:  ip2,
				}
			}
		}
	} else {
		var count *big.Int
		if !strings.Contains(s, "/") {
			if !strings.Contains(s, ":") {
				s = s + "/32"
			} else {
				s = s + "/128"
			}
			count = big.NewInt(1)
		}
		ip, ipNet, err := net.ParseCIDR(s)
		if err == nil && !strings.HasSuffix(s, "/0") {
			return &IpAddr{
				Kind:  IpaddrKindCidr,
				Ip:    ip,
				IpNet: ipNet,
				count: count,
			}
		}
	}
	return nil
}

func (o *IpAddr) Count() *big.Int {
	if o.count != nil {
		return o.count
	}

	result := big.NewInt(0)

	if o.Kind == IpaddrKindCidr {
		flag := false
		for _, m := range o.Mask() {
			result = result.Lsh(result, 8)
			if m != 255 && !flag {
				result = result.Add(result, big.NewInt(256-int64(m))) //uint64(o.Ip[i]|^m) - uint64(o.Ip[i]&m) + 1
				flag = true
			}
		}
		if !flag {
			result = big.NewInt(1)
		}
	} else {
		for i := range o.Ip1 {
			if o.Ip1[i] < o.Ip2[i] {
				result = result.Add(result, big.NewInt(int64(o.Ip2[i]-o.Ip1[i]-1)))
				for j := i + 1; j < len(o.Ip1); j++ {
					result = result.Lsh(result, 8)
					result = result.Add(result, big.NewInt(int64(255-o.Ip1[j])+int64(o.Ip2[j]-0)))
				}
				result = result.Add(result, big.NewInt(2))
				break
			}
		}
	}
	o.count = result
	return result
}

func (o *IpAddr) Mask() []byte {
	// assert(o.Kind == IpaddrKindCidr)
	if len(o.IpNet.Mask) == 16 {
		return o.IpNet.Mask
	}
	return append([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, o.IpNet.Mask...)
}

func (o *IpAddr) String() string {
	if o.Kind == IpaddrKindRange {
		return o.Ip1.String() + "-" + o.Ip2.String()
	} else if o.Count().Cmp(big.NewInt(1)) == 0 {
		return o.Ip.String()
	} else {
		return o.IpNet.String()
	}
}

func compareIP(ip1, ip2 net.IP) int {
	for i := 0; i < len(ip1); i++ {
		if ip1[i] > ip2[i] {
			return 1
		} else if ip1[i] < ip2[i] {
			return -1
		}
	}
	return 0
}
