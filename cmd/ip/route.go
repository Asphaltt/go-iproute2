package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/Asphaltt/go-iproute2/ip"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

func init() {
	rootCmd.AddCommand(routeCmd())
}

func routeCmd() *cobra.Command {
	routeCmd := &cobra.Command{
		Use:     "route",
		Aliases: []string{"r", "ro", "rou", "rout"},
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listRoutes)
		},
	}
	routeCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "li", "lis", "lst", "s", "sh", "sho", "show"},
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listRoutes)
		},
	})
	return routeCmd
}

func (c *client) listRoutes() {
	ipcli := ip.NewWithConn(c.conn)
	entries, err := ipcli.ListRoutes()
	if err != nil {
		fmt.Println("failed to list route entries, err:", err)
		return
	}

	for _, e := range entries {
		printRouteEntry(e)
	}
}

func printRouteEntry(e *ip.RouteEntry) {
	table := e.Table
	if table == -1 {
		table = e.TableID
	}
	// if int(table) != unix.RT_TABLE_MAIN {
	// 	continue
	// }

	var s strings.Builder

	if e.Type != unix.RTN_UNICAST {
		s.WriteString(fmt.Sprintf("%s ", e.Type))
	}

	if e.Daddr != nil {
		if e.DstLen != 32 { // IPv4 地址有 32 个比特
			s.WriteString(fmt.Sprintf("%s/%d ", e.Daddr, e.DstLen))
		} else {
			s.WriteString(fmt.Sprintf("%s ", e.Daddr))
		}
	} else if e.DstLen != 0 {
		s.WriteString(fmt.Sprintf("0/%d ", e.DstLen))
	} else {
		s.WriteString("default ")
	}

	if e.Saddr != nil {
		if e.SrcLen != 32 {
			s.WriteString(fmt.Sprintf("from %s/%d ", e.Saddr, e.SrcLen))
		} else {
			s.WriteString(fmt.Sprintf("from %s ", e.Saddr))
		}
	} else if e.SrcLen != 0 {
		s.WriteString(fmt.Sprintf("from 0/%d ", e.SrcLen))
	}

	if e.Tos != 0 {
		s.WriteString(fmt.Sprintf("tos %d ", e.Tos))
	}

	if e.Gateway != nil {
		s.WriteString(fmt.Sprintf("via %s ", e.Gateway))
	}

	if e.OutIfindex != 0 {
		ifi, err := net.InterfaceByIndex(e.OutIfindex)
		if err == nil {
			s.WriteString(fmt.Sprintf("dev %s ", ifi.Name))
		}
	}

	if int(table) != unix.RT_TABLE_MAIN {
		s.WriteString(fmt.Sprintf("table %s ", table))
	}

	if int(e.Flags)&unix.RTM_F_CLONED == 0 {
		if int(e.Protocol) != unix.RTPROT_BOOT {
			s.WriteString(fmt.Sprintf("proto %s ", e.Protocol))
		}

		if int(e.Scope) != unix.RT_SCOPE_UNIVERSE {
			s.WriteString(fmt.Sprintf("scope %s ", e.Scope))
		}
	}

	if e.PrefSrc != nil {
		s.WriteString(fmt.Sprintf("src %s ", e.PrefSrc))
	}

	if e.Priority != -1 {
		s.WriteString(fmt.Sprintf("metric %d ", e.Priority))
	}

	s.WriteString(e.Flags.String())

	if e.Pref != -1 {
		s.WriteString(fmt.Sprintf("pref %s ", e.Pref))
	}

	fmt.Println(s.String())
}
