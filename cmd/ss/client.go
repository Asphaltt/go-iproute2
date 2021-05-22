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

	mix   bool
	netid string
}

func dialNetlink() (*client, error) {
	conn, err := netlink.Dial(iproute2.FamilySocketMonitoring, nil)
	if err != nil {
		return nil, err
	}

	var c client
	c.conn = conn
	c.mix = isMixConfig()
	return &c, nil
}

func (c *client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *client) showSocketInfoHeader() {
	if c.mix {
		// "Netid     "
		fmt.Printf("Netid     ")
	}
	// "State     Recv-Q     Send-Q     Local Address:Port     Peer Address:Port"
	fmt.Printf("%-10s     %-6s     %-6s    %24s:%-5s     %24s:%-5s\n",
		"State", "Recv-Q", "Send-Q",
		"Local Address", "Port", "Peer Address", "Port")
}

func (c *client) showEntries(entries []*ss.Entry, err error) error {
	if err != nil {
		return err
	}
	for _, e := range entries {
		printEntry(e, c.netid)
	}
	return nil
}

// "State     Recv-Q     Send-Q     Local Address:Port     Peer Address:Port"
const fmtEntry = "%-10s     %-6d     %-6d    %24s:%-5d     %24s:%s\n"

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
func printEntry(e *ss.Entry, netid string) {
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
		peerPort = "*"
	} else {
		peerPort = fmt.Sprintf("%-5d", e.PeerPort)
	}
	fmt.Printf(fmtEntry,
		e.State.String(),
		e.RecvQ, e.SendQ,
		laddr, e.LocalPort,
		peerAddr, peerPort,
	)
}
