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

var rootCmd = cobra.Command{
	Use: "bridge",
}

func main() {
	rootCmd.Execute()
}
