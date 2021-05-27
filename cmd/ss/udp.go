package main

import (
	"fmt"

	"github.com/Asphaltt/go-iproute2/ss"
)

func (c *client) showUDP(listen, udp4, udp6 bool) {
	if listen {
		c.showUDPListeners(udp4, udp6)
	} else {
		c.showUDPSockets(udp4, udp6)
	}
}

func (c *client) showUDPSockets(udp4, udp6 bool) {
	sc := ss.New(c.conn)
	if udp4 {
		c.showUDP4Sockets(sc)
	}
	if udp6 {
		c.showUDP6Sockets(sc)
	}
}

func (c *client) showUDP4Sockets(sc *ss.Client) {
	if err := c.showEntries(sc.ListUdp4Sockets()); err != nil {
		fmt.Println("failed to list IPv4 udp sockets, err:", err)
	}
}

func (c *client) showUDP6Sockets(sc *ss.Client) {
	if err := c.showEntries(sc.ListUdp6Sockets()); err != nil {
		fmt.Println("failed to list IPv6 udp sockets, err:", err)
	}
}

func (c *client) showUDPListeners(udp4, udp6 bool) {
	if c.mix {
		c.netid = "udp"
	}
	sc := ss.New(c.conn)
	if udp4 {
		c.showUDP4Listeners(sc)
	}
	if udp6 {
		c.showUDP6Listeners(sc)
	}
}

func (c *client) showUDP4Listeners(sc *ss.Client) {
	if err := c.showEntries(sc.ListUdp4Listeners()); err != nil {
		fmt.Println("failed to list IPv4 udp listeners, err:", err)
	}
}

func (c *client) showUDP6Listeners(sc *ss.Client) {
	if err := c.showEntries(sc.ListUdp6Listeners()); err != nil {
		fmt.Println("failed to list IPv6 udp listeners, err:", err)
	}
}
