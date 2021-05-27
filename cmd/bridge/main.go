package main

import (
	"fmt"

	"github.com/Asphaltt/go-iproute2"
	"github.com/Asphaltt/go-iproute2/bridge"
	"github.com/mdlayher/netlink"
	"github.com/spf13/cobra"
)

var cli client

type client struct {
	conn *netlink.Conn
}

func (c *client) dialNetlink() error {
	var err error
	c.conn, err = netlink.Dial(iproute2.NETLINK_ROUTE, nil)
	return err
}

func (c *client) dialFdbMonitor() error {
	var err error
	c.conn, err = bridge.DialFdbMonitor()
	return err
}

func (c *client) runCmd(fn func())        { c.run(c.dialNetlink, fn) }
func (c *client) runFdbMonitor(fn func()) { c.run(c.dialFdbMonitor, fn) }
func (c *client) run(dial func() error, fn func()) {
	err := dial()
	if err != nil {
		fmt.Println("failed to create netlink socket, err:", err)
		return
	}
	defer c.conn.Close()

	fn()
}

func main() {
	rootCmd := &cobra.Command{
		Use: "bridge",
	}

	monitorCmd := &cobra.Command{
		Use: "monitor",
	}
	monitorCmd.AddCommand(&cobra.Command{
		Use: "fdb",
		Run: func(cmd *cobra.Command, args []string) {
			cli.runFdbMonitor(cli.monitorFdb)
		},
	})

	fdbCmd := &cobra.Command{
		Use: "fdb",
	}
	fdbCmd.AddCommand(&cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listFdb)
		},
	})

	rootCmd.AddCommand(fdbCmd)
	rootCmd.AddCommand(monitorCmd)
	rootCmd.Execute()
}
