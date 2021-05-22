package main

import (
	"fmt"
	"os"

	"github.com/Asphaltt/go-iproute2"
	"github.com/mdlayher/netlink"
)

type client struct {
	conn *netlink.Conn
}

func dialNetlink() (*client, error) {
	conn, err := netlink.Dial(iproute2.FamilySocketMonitoring, nil)
	if err != nil {
		return nil, err
	}

	return &client{conn}, nil
}

func (c *client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

var config struct {
	v4, v6 bool
	listen bool
	tcp    bool
}

func main() {
	config.v4 = true
	config.v6 = true

	for _, arg := range os.Args[1:] {
		if arg[0] != '-' {
			continue
		}

		for _, b := range arg[1:] {
			switch b {
			case '4':
				config.v4 = true
				config.v6 = false
			case '6':
				config.v4 = false
				config.v6 = true
			case 'l':
				config.listen = true
			case 't':
				config.tcp = true
			}
		}
	}

	c, err := dialNetlink()
	if err != nil {
		fmt.Println("failed to create netlink socket, err:", err)
		return
	}
	defer c.Close()

	showSocketInfoHeader()
	if config.tcp {
		if config.listen {
			c.showTCPListeners(config.v4, config.v6)
		} else {
			c.showTCPConns(config.v4, config.v6)
		}
	}
}

func showSocketInfoHeader() {
	// "State     Recv-Q     Send-Q     Local Address:Port     Peer Address:Port"
	fmt.Printf("%-10s     %-6s     %-6s    %24s:%-5s     %24s:%-5s\n",
		"State", "Recv-Q", "Send-Q",
		"Local Address", "Port", "Peer Address", "Port")
}
