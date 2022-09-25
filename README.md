# iptrie
a fast ip matcher using trie

### Supporting ipv4/ipv6 range, cidr 
```
"192.168.100.2"
"192.168.10.10-192.168.20.9"
"172.16.100.1/24"
"2001:db8::68"
```

### Example
```
package main

import "github.com/leetjob/iptrie"

func main() {
	ipList := []string{
		"192.168.100.2",
		"192.168.10.10-192.168.20.9",
		"172.16.100.0/24",
		"2001:db8::68",
	}

	ipTrie := NewTrie()
	for _, ip := range ipList {
		ipTrie.InsertIpAddr(NewIpAddr(ip), nil)
	}
	result := ipTrie.FindByIp(net.ParseIP("172.16.100.100"))
	fmt.Printf("result=%v", result)
}
```