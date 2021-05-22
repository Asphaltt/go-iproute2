package main

import (
	"fmt"

	"github.com/Asphaltt/go-iproute2/ss"
)

func (c *client) showTCPConns(tcp4, tcp6 bool) {
	sc := ss.NewClient(c.conn)
	if tcp4 {
		c.showTCP4Conns(sc)
	}
	if tcp6 {
		c.showTCP6Conns(sc)
	}
}

func (c *client) showTCP4Conns(sc *ss.Client) {
	if err := c.showTCPEntries(sc.ListTcp4Conns()); err != nil {
		fmt.Println("failed to list IPv4 tcp connections, err:", err)
	}
}

func (c *client) showTCP6Conns(sc *ss.Client) {
	if err := c.showTCPEntries(sc.ListTcp6Conns()); err != nil {
		fmt.Println("failed to list IPv6 tcp connections, err:", err)
	}
}

func (c *client) showTCPListeners(tcp4, tcp6 bool) {
	sc := ss.NewClient(c.conn)
	if tcp4 {
		c.showTCP4Listeners(sc)
	}
	if tcp6 {
		c.showTCP6Listeners(sc)
	}
}

func (c *client) showTCP4Listeners(sc *ss.Client) {
	if err := c.showTCPEntries(sc.ListTcp4Listeners()); err != nil {
		fmt.Println("failed to list IPv4 tcp listeners, err:", err)
	}
}

func (c *client) showTCP6Listeners(sc *ss.Client) {
	if err := c.showTCPEntries(sc.ListTcp6Listeners()); err != nil {
		fmt.Println("failed to list IPv6 tcp listeners, err:", err)
	}
}

func (c *client) showTCPEntries(entries []*ss.TcpEntry, err error) error {
	if err != nil {
		return err
	}
	for _, e := range entries {
		fmt.Println(e.String())
	}
	return nil
}
