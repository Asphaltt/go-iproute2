package main

import (
	"fmt"
	"os"

	"github.com/Asphaltt/go-iproute2"
	"github.com/Asphaltt/go-iproute2/ss"
	"github.com/mdlayher/netlink"
)

var config struct {
	tcp bool
}

func main() {
	for _, arg := range os.Args[1:] {
		if arg[0] != '-' {
			continue
		}

		for _, b := range arg[1:] {
			switch b {
			case 't':
				config.tcp = true
			}
		}
	}

	showSocketInfoHeader()
	if config.tcp {
		showTCP()
	}
}

func showTCP() {
	conn, err := netlink.Dial(iproute2.FamilySocketMonitoring, nil)
	if err != nil {
		fmt.Println("failed to dial a socket monitoring netlink connection, err:", err)
		return
	}
	defer conn.Close()

	entries, err := ss.ListTcpConns(conn)
	if err != nil {
		fmt.Println("failed to list tcp connections, err:", err)
		return
	}

	for _, e := range entries {
		fmt.Println(e.String())
	}
}

func showSocketInfoHeader() {
	// "State     Recv-Q     Send-Q     Local Address:Port     Peer Address:Port"
	fmt.Printf("%-10s     %-6s     %-6s    %24s:%-5s     %24s:%-5s\n",
		"State", "Recv-Q", "Send-Q",
		"Local Address", "Port", "Peer Address", "Port")
}
