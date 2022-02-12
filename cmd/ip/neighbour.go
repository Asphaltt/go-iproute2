package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/Asphaltt/go-iproute2"
	"github.com/Asphaltt/go-iproute2/ip"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(neighCmd())
}

func neighCmd() *cobra.Command {
	neighCmd := &cobra.Command{
		Use:     "neighbour",
		Aliases: []string{"n", "ne", "nei", "neig", "neigh"},
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listNeighbours)
		},
	}
	neighCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "li", "lis", "lst", "s", "sh", "sho", "show"},
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listNeighbours)
		},
	})
	return neighCmd
}

func (c *client) listNeighbours() {
	ipcli := ip.NewWithConn(c.conn)
	entries, err := ipcli.ListNeighbours()
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
