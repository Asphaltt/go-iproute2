package main

import (
	"fmt"

	"github.com/Asphaltt/go-iproute2/ss"
)

func (c *client) showTCPConns() {
	ssCli := ss.NewClient(c.conn)
	entries, err := ssCli.ListTcpConns()
	if err != nil {
		fmt.Println("failed to list tcp connections, err:", err)
		return
	}

	for _, e := range entries {
		fmt.Println(e.String())
	}
}

func (c *client) showTCPListeners() {
	ssCli := ss.NewClient(c.conn)
	entries, err := ssCli.ListTcp4Listeners()
	if err != nil {
		fmt.Println("failed to list IPv4 tcp listeners, err:", err)
		return
	}
	for _, e := range entries {
		fmt.Println(e.String())
	}

	entries, err = ssCli.ListTcp6Listeners()
	if err != nil {
		fmt.Println("failed to list IPv6 tcp listeners, err:", err)
		return
	}
	for _, e := range entries {
		fmt.Println(e.String())
	}
}
