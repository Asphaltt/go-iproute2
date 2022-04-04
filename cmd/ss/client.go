package main

import (
	"fmt"
	"net"

	"github.com/Asphaltt/go-iproute2"
	"github.com/Asphaltt/go-iproute2/ss"
	"github.com/mdlayher/netlink"
)

type client struct {
	conn *netlink.Conn
}

func dialNetlink() (*client, error) {
	conn, err := netlink.Dial(iproute2.FamilySocketMonitoring, nil)
	if err != nil {
		return nil, err
	}

	var c client
	c.conn = conn
	return &c, nil
}

func (c *client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *client) showSocketInfoHeader() {
	if isMixConfig() {
		// "Netid     "
		fmt.Printf("Netid     ")
	}
	// "State     Recv-Q     Send-Q     Local Address:Port     Peer Address:Port"
	fmt.Printf("%-10s     %-6s     %-6s    %24s:%-5s     %24s:%-5s",
		"State", "Recv-Q", "Send-Q",
		"Local Address", "Port", "Peer Address", "Port")

	// show process info
	if config.process {
		fmt.Print(" Process")
	}

	fmt.Println()
}

func getNetid(netid string) string {
	if isMixConfig() {
		return netid
	}
	return ""
}

func (c *client) showTcpEntries(fn func() ([]*ss.Entry, error), showProc bool) error {
	return c.showEntries(fn, getNetid("tcp"), showProc)
}

func (c *client) showUdpEntries(fn func() ([]*ss.Entry, error), showProc bool) error {
	return c.showEntries(fn, getNetid("udp"), showProc)
}

func (c *client) showEntries(fn func() ([]*ss.Entry, error), netid string, showProc bool) error {
	entries, err := fn()
	if err != nil {
		return err
	}
	for _, e := range entries {
		printEntry(e, netid, showProc)
	}
	return nil
}

// "State     Recv-Q     Send-Q     Local Address:Port     Peer Address:Port"
const fmtEntry = "%-10s     %-6d     %-6d    %24s:%-5d     %24s:%s"

func isZeroIPv6Addr(addr net.IP) bool {
	if addr.To16() == nil {
		return false
	}

	for i := 0; i < len(addr); i++ {
		if addr[i] != 0 {
			return false
		}
	}
	return true
}

// String formats Entry into
// '[Netid    ]State     Recv-Q     Send-Q     Local Address:Port     Peer Address:Port'.
func printEntry(e *ss.Entry, netid string, showProc bool) {
	if netid != "" {
		fmt.Printf("%-5s     ", netid)
	}
	var laddr string
	if e.Ifindex == 0 {
		if !e.IsIPv4 {
			laddr = fmt.Sprintf("[%s]", e.LocalAddr.To16().String())
		} else {
			laddr = e.LocalAddr.String()
		}
	} else {
		ifc, _ := net.InterfaceByIndex(e.Ifindex)
		if !e.IsIPv4 {
			if isZeroIPv6Addr(e.LocalAddr) {
				laddr = fmt.Sprintf("*%%%s", ifc.Name)
			} else {
				laddr = fmt.Sprintf("[%s]%%%s",
					e.LocalAddr.To16().String(), ifc.Name)
			}
		} else {
			laddr = fmt.Sprintf("%s%%%s", e.LocalAddr.String(), ifc.Name)
		}
	}
	var peerAddr string
	if !e.IsIPv4 {
		peerAddr = fmt.Sprintf("[%s]", e.PeerAddr.To16().String())
	} else {
		peerAddr = e.PeerAddr.String()
	}
	var peerPort string
	if e.PeerPort == 0 {
		peerPort = fmt.Sprintf("%-5s", "*")
	} else {
		peerPort = fmt.Sprintf("%-5d", e.PeerPort)
	}
	fmt.Printf(fmtEntry,
		e.State.String(),
		e.RecvQ, e.SendQ,
		laddr, e.LocalPort,
		peerAddr, peerPort,
	)

	// show process info
	if showProc {
		if pinfo, ok := procs[e.Inode]; ok {
			fmt.Printf(" users:(%s)", pinfo)
		}
	}

	fmt.Println() // newline
}
