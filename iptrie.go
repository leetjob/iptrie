package iptrie

import "net"

type IpTrieNode struct {
	mark bool
	data interface{}
	next []*IpTrieNode
}

type IpTrie struct {
	root     *IpTrieNode
	fastRoot *IpTrieNode
}

func NewTrie() *IpTrie {
	t := &IpTrie{root: &IpTrieNode{next: make([]*IpTrieNode, 256)}}
	t.setFastRoot([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0x7f, 0x0, 0x0, 0x1}, 11)
	return t
}

func (o *IpTrie) setFastRoot(keyword []byte, k int) {
	cur := o.root
	for i, b := range keyword {
		if cur.next[b] == nil {
			cur.next[b] = &IpTrieNode{next: make([]*IpTrieNode, 256)}
		}
		cur = cur.next[b]
		if i == k {
			o.fastRoot = cur
		}
	}
}

func (o *IpTrie) InsertIpAddr(ipAddr *IpAddr, data interface{}) {
	if ipAddr.Kind == IpaddrKindCidr {
		o.insert(ipAddr.Ip, ipAddr.Mask(), data)
	} else {
		o.insertRange(ipAddr.Ip1, ipAddr.Ip2, data)
	}
}

func (o *IpTrie) FindByIp(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if v4 := ip.To4(); v4 != nil {
		return o.find(o.fastRoot, v4)
	} else {
		return o.find(o.root, ip)
	}
}

func (o *IpTrie) FindLowestByIp(ip net.IP) interface{} {
	if ip == nil {
		return nil
	}
	if v4 := ip.To4(); v4 != nil {
		return o.findLowest(o.fastRoot, v4)
	} else {
		return o.findLowest(o.root, ip)
	}
}

func (o *IpTrie) insert(keyword []byte, mask []byte, data interface{}) {
	cur := o.root
	for i, b := range keyword {
		if cur.mark {
			return
		}

		if mask[i] == 255 {
			if cur.next[b] == nil {
				cur.next[b] = &IpTrieNode{next: make([]*IpTrieNode, 256)}
			}
			cur = cur.next[b]
		} else {
			// [b&mask[i], b|^mask[i]] mark true
			if mask[i] == 0 {
				cur.mark, cur.data = true, data
			} else {
				for _, mb := range slice(b&mask[i], b|^mask[i]) {
					if cur.next[mb] == nil {
						cur.next[mb] = &IpTrieNode{next: make([]*IpTrieNode, 256)}
					}
					if !cur.next[mb].mark {
						cur.next[mb].mark, cur.next[mb].data = true, data
					}
				}
			}
			return
		}
	}
	if !cur.mark {
		cur.mark, cur.data = true, data
	}
}

func (o *IpTrie) insertRange(keyword1 []byte, keyword2 []byte, data interface{}) {
	cur := o.root
	for i, b1 := range keyword1 {
		b2 := keyword2[i]

		if b1 == b2 {
			if cur.next[b1] == nil {
				cur.next[b1] = &IpTrieNode{next: make([]*IpTrieNode, 256)}
			}
			cur = cur.next[b1]
			if cur.mark {
				return
			}
		} else if b1 < b2 {
			//(b1, b2) mark true
			for _, b := range slice(b1, b2) {
				if cur.next[b] == nil {
					cur.next[b] = &IpTrieNode{next: make([]*IpTrieNode, 256)}
				}
				if b > b1 && b < b2 && !cur.next[b].mark {
					cur.next[b].mark, cur.next[b].data = true, data
				}
			}

			// keyword1
			cur1 := cur.next[b1]
			for j := i + 1; j < len(keyword1); j++ {
				if cur1.mark {
					break
				}
				for _, b := range slice(keyword1[j], 255) {
					if cur1.next[b] == nil {
						cur1.next[b] = &IpTrieNode{next: make([]*IpTrieNode, 256)}
					}
					if b > keyword1[j] && !cur1.next[b].mark {
						cur1.next[b].mark, cur1.next[b].data = true, data
					}
				}
				cur1 = cur1.next[keyword1[j]]
			}
			if !cur1.mark {
				cur1.mark, cur1.data = true, data
			}

			// keyword2
			cur2 := cur.next[b2]
			for j := i + 1; j < len(keyword2); j++ {
				if cur2.mark {
					break
				}
				for _, b := range slice(0, keyword2[j]) {
					if cur2.next[b] == nil {
						cur2.next[b] = &IpTrieNode{next: make([]*IpTrieNode, 256)}
					}
					if b < keyword2[j] && !cur2.next[b].mark {
						cur2.next[b].mark, cur2.next[b].data = true, data
					}
				}
				cur2 = cur2.next[keyword2[j]]
			}

			if !cur2.mark {
				cur2.mark, cur2.data = true, data
			}

			break
		} else {
			//never get here
		}
	}
}

func (o *IpTrie) find(root *IpTrieNode, keyword []byte) bool {
	cur := root
	for _, b := range keyword {
		if cur.next[b] == nil {
			return false
		}
		cur = cur.next[b]
		if cur.mark {
			return true
		}
	}
	return false
}

func (o *IpTrie) findLowest(root *IpTrieNode, keyword []byte) interface{} {
	var data interface{}
	cur := root
	for _, b := range keyword {
		if cur.next[b] == nil {
			return data
		}
		cur = cur.next[b]
		if cur.mark {
			data = cur.data
		}
	}
	return data
}

// [low, high]
func slice(low, high byte) []byte {
	s := make([]byte, 0)
	if low > high {
		return s
	}
	for i := low; i != high; i++ {
		s = append(s, i)
	}
	return append(s, high)
}
