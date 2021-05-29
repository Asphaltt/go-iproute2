package main

import (
	"fmt"

	"github.com/Asphaltt/go-iproute2"
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

func (c *client) runCmd(fn func()) {
	err := c.dialNetlink()
	if err != nil {
		fmt.Println("failed to create netlink socket, err:", err)
		return
	}
	defer c.conn.Close()

	fn()
}

func main() {
	rootCmd := &cobra.Command{
		Use: "ip",
	}

	neighCmd := &cobra.Command{
		Use: "neigh",
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listNeighbours)
		},
	}
	neighCmd.AddCommand(&cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listNeighbours)
		},
	})

	linkCmd := &cobra.Command{
		Use: "link",
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listLinks)
		},
	}
	linkCmd.AddCommand(&cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listLinks)
		},
	})

	rootCmd.AddCommand(neighCmd)
	rootCmd.AddCommand(linkCmd)
	rootCmd.Execute()
}
