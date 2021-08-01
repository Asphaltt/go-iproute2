package main

import (
	"fmt"
	"net"

	"github.com/Asphaltt/go-iproute2/bridge"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(monitorCmd())
}

func monitorCmd() *cobra.Command {
	monitorCmd := &cobra.Command{
		Use: "monitor",
	}
	monitorCmd.AddCommand(&cobra.Command{
		Use: "fdb",
		Run: func(cmd *cobra.Command, args []string) {
			cli.runFdbMonitor(cli.monitorFdb)
		},
	})
	return monitorCmd
}

func (c *client) monitorFdb() {
	bcli := bridge.New(c.conn)
	err := bcli.MonitorFdb(printFdbEntry)
	if err != nil {
		fmt.Println("failed to bridge monitor fdb, err:", err)
	}
}

func printFdbEntry(entry *bridge.FdbEntry) {
	var action string
	switch entry.Action {
	case bridge.FdbActionAdd:
		action = "Added"
	case bridge.FdbActionDel:
		action = "Deleted"
	default:
		action = "Unkowned"
	}
	devInfo, _ := net.InterfaceByIndex(entry.Ifindex)
	masterInfo, _ := net.InterfaceByIndex(entry.Master)
	fmt.Printf("%s %s dev %s master %s\n", action, entry.Lladdr, devInfo.Name, masterInfo.Name)
}
