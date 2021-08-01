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
	if ifc, err := net.InterfaceByIndex(e.Ifindex); err == nil {
		devName = ifc.Name
	}

	// TODO: show vlan
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s dev %s", e.Lladdr, devName))
	if e.Flag != 0 {
		b.WriteString(fmt.Sprintf(" %s", e.Flag))
	}
	b.WriteString(fmt.Sprintf(" %s", e.State))
	fmt.Println(b.String())
}
