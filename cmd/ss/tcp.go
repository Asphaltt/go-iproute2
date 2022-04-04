package main

import (
	"fmt"

	"github.com/Asphaltt/go-iproute2/ss"
)

func (c *client) showTCP(listen, tcp4, tcp6 bool) {
	if listen {
		c.showTCPListeners(tcp4, tcp6)
	} else {
		c.showTCPConns(tcp4, tcp6)
	}
}

func (c *client) showTCPConns(tcp4, tcp6 bool) {
	sc := ss.NewWithConn(c.conn)
	if tcp4 {
		c.showTCP4Conns(sc)
	}
	if tcp6 {
		c.showTCP6Conns(sc)
	}
}

func (c *client) showTCP4Conns(sc *ss.Client) {
	if err := c.showTcpEntries(sc.ListTcp4Conns, config.process); err != nil {
		fmt.Println("failed to list IPv4 tcp connections, err:", err)
	}
}

func (c *client) showTCP6Conns(sc *ss.Client) {
	if err := c.showTcpEntries(sc.ListTcp6Conns, config.process); err != nil {
		fmt.Println("failed to list IPv6 tcp connections, err:", err)
	}
}

func (c *client) showTCPListeners(tcp4, tcp6 bool) {
	sc := ss.NewWithConn(c.conn)
	if tcp4 {
		c.showTCP4Listeners(sc)
	}
	if tcp6 {
		c.showTCP6Listeners(sc)
	}
}

func (c *client) showTCP4Listeners(sc *ss.Client) {
	if err := c.showTcpEntries(sc.ListTcp4Listeners, config.process); err != nil {
		fmt.Println("failed to list IPv4 tcp listeners, err:", err)
	}
}

func (c *client) showTCP6Listeners(sc *ss.Client) {
	if err := c.showTcpEntries(sc.ListTcp6Listeners, config.process); err != nil {
		fmt.Println("failed to list IPv6 tcp listeners, err:", err)
	}
}
