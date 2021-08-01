package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/Asphaltt/go-iproute2/bridge"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fdbCmd())
}

func fdbCmd() *cobra.Command {
	fdbCmd := &cobra.Command{
		Use: "fdb",
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listFdb)
		},
	}
	fdbCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls", "lis", "lst", "s", "sh", "sho", "show"},
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listFdb)
		},
	})
	return fdbCmd
}

func (c *client) listFdb() {
	bcli := bridge.New(c.conn)
	entries, err := bcli.ListFdb()
	if err != nil {
		fmt.Println("failed to list fdb entries, err:", err)
		return
	}

	for _, e := range entries {
		printListFdb(e)
	}
}

func printListFdb(e *bridge.FdbEntry) {
	var devName string
	if ifi, err := net.InterfaceByIndex(e.Ifindex); err == nil {
		devName = ifi.Name
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s dev %s ", e.Lladdr, devName))

	if e.Vlan != 0 {
		b.WriteString(fmt.Sprintf("vlan %d ", e.Vlan))
	}

	if e.Flag != 0 {
		b.WriteString(fmt.Sprintf("%s ", e.Flag))
	}

	if e.Master != 0 {
		ifi, err := net.InterfaceByIndex(e.Master)
		if err == nil {
			b.WriteString(fmt.Sprintf("master %s ", ifi.Name))
		}
	}

	b.WriteString(e.State.String())

	fmt.Println(b.String())
}
