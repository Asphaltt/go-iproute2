package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/Asphaltt/go-iproute2"
	"github.com/Asphaltt/go-iproute2/ip"
)

func (c *client) listNeighbours() {
	ipCli := ip.New(c.conn)
	entries, err := ipCli.ListNeighbours()
	if err != nil {
		fmt.Println("failed to list neighbour entries, err:", err)
		return
	}

	for _, e := range entries {
		if e.State != iproute2.NudNoArp {
			printNeighEntry(e)
		}
	}
}

func printNeighEntry(e *ip.NeighEntry) {
	var b strings.Builder
	b.WriteString(e.Addr.String())

	if e.Ifindex != 0 {
		ifc, _ := net.InterfaceByIndex(e.Ifindex)
		b.WriteString(fmt.Sprintf(" dev %s", ifc.Name))
	}

	if e.State != iproute2.NudFailed {
		b.WriteString(fmt.Sprintf(" lladdr %s", e.Lladdr))
	}
	b.WriteString(fmt.Sprintf(" %s", strings.ToUpper(e.State.String())))
	fmt.Println(b.String())
}
