package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/Asphaltt/go-iproute2/ip"
	"golang.org/x/sys/unix"
)

func (c *client) listLinks() {
	ipcli := ip.New(c.conn)
	entries, err := ipcli.ListLinks()
	if err != nil {
		fmt.Println("failed to list link entries, err:", err)
		return
	}

	for _, e := range entries {
		printLinkEntry(e)
	}
}

func printLinkEntry(e *ip.LinkEntry) {
	if e.Name == "" {
		return
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("%d: ", e.Ifindex))
	if e.Link != 0 {
		if e.Namespace >= 0 {
			s.WriteString(fmt.Sprintf("%s@if%d: ", e.Name, e.Link))
		} else {
			ifinfo, err := net.InterfaceByIndex(e.Link)
			if err == nil {
				s.WriteString(fmt.Sprintf("%s@%s: ", e.Name, ifinfo.Name))
			}
		}
	} else {
		s.WriteString(fmt.Sprintf("%s: ", e.Name))
	}
	s.WriteString(fmt.Sprintf("%s ", e.DeviceFlags))
	if e.MTU != 0 {
		s.WriteString(fmt.Sprintf("mtu %d ", e.MTU))
	}
	if e.QDisc != "" {
		s.WriteString(fmt.Sprintf("qdisc %s ", e.QDisc))
	}
	if e.Master != 0 {
		ifinfo, err := net.InterfaceByIndex(e.Master)
		if err == nil {
			s.WriteString(fmt.Sprintf("master %s ", ifinfo.Name))
		}
	}
	if e.OperState >= 0 {
		s.WriteString(fmt.Sprintf("state %s ", e.OperState))
	}
	s.WriteString(fmt.Sprintf("mode %s ", e.Mode))
	if e.Group >= 0 {
		s.WriteString(fmt.Sprintf("group %s ", e.Group))
	}
	if e.TxQueue != 0 {
		s.WriteString(fmt.Sprintf("qlen %d", e.TxQueue))
	}

	s.WriteString("\n")
	s.WriteString(fmt.Sprintf("    link/%s ", e.DeviceType))
	if e.Addr != nil {
		switch len(e.Addr) {
		case 4, 16:
			s.WriteString(net.IP(e.Addr).String())
		default:
			s.WriteString(net.HardwareAddr(e.Addr).String())
		}
	}
	if e.Broadcast != nil {
		if e.DeviceFlags&unix.IFF_POINTOPOINT != 0 {
			s.WriteString(" peer link_pointtopoint")
		} else {
			switch len(e.Broadcast) {
			case 4, 16:
				s.WriteString(fmt.Sprintf(" brd %s", net.IP(e.Broadcast)))
			default:
				s.WriteString(fmt.Sprintf(" brd %s", net.HardwareAddr(e.Broadcast)))
			}
		}
	}
	if e.Namespace >= 0 {
		s.WriteString(fmt.Sprintf(" link-netnsid %d", e.Namespace))
	}
	fmt.Println(s.String())
}
