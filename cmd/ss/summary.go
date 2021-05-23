package main

import (
	"fmt"

	"github.com/Asphaltt/go-iproute2/ss"
)

func showSummary() {
	var s ss.Summary
	if err := s.GetSockStat(); err != nil {
		fmt.Println("failed to get sock stat, err:", err)
		return
	}
	if err := s.GetSockStat6(); err != nil {
		fmt.Println("failed to get sock6 stat, err:", err)
		return
	}
	if err := s.GetTcpEstablished(); err != nil {
		fmt.Println("failed to get number of established tcp connections, err:", err)
		return
	}
	printSummary(&s)
}

/*
 * summary example:
 *
 * Total: 182
 * TCP:   16 (estab 3, closed 1, orphaned 0, synrecv 0, timewait 0/0), ports 0
 *
 * Transport Total     IP        IPv6
 * RAW       0         0         0
 * UDP       11        8         3
 * TCP       15        12        3
 * INET      26        20        6
 * FRAG      0         0         0
 *
 */
func printSummary(s *ss.Summary) {
	fmt.Printf("Total: %d\n", s.Sockets)
	fmt.Printf("TCP:   %d (estab %d, closed %d, orphaned %d, timewait %d)\n",
		s.TcpTotal+s.TcpTws, s.TcpEstablished,
		s.TcpTotal-(s.Tcp4Hashed+s.Tcp6Hashed-s.TcpTws),
		s.TcpOrphans, s.TcpTws)
	fmt.Println()

	fmt.Println("Transport Total     IP        IPv6")
	fmt.Printf("RAW	  %-9d %-9d %-9d\n", s.Raw4+s.Raw6, s.Raw4, s.Raw6)
	fmt.Printf("UDP	  %-9d %-9d %-9d\n", s.Udp4+s.Udp6, s.Udp4, s.Udp6)
	fmt.Printf("TCP	  %-9d %-9d %-9d\n",
		s.Tcp4Hashed+s.Tcp6Hashed, s.Tcp4Hashed, s.Tcp6Hashed)
	fmt.Printf("INET	  %-9d %-9d %-9d\n",
		s.Raw4+s.Udp4+s.Tcp4Hashed+s.Raw6+s.Udp6+s.Tcp6Hashed,
		s.Raw4+s.Udp4+s.Tcp4Hashed,
		s.Raw6+s.Udp6+s.Tcp6Hashed)
	fmt.Printf("FRAG	  %-9d %-9d %-9d\n", s.Frag4+s.Frag6, s.Frag4, s.Frag6)
	fmt.Println()
}
