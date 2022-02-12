package main

import (
	"fmt"
	"sort"
	"strings"
	"syscall"

	"github.com/Asphaltt/go-iproute2/ip"
	"github.com/spf13/cobra"
)

const (
	INFINITY_LIFE_TIME = 0xFFFFFFFF
)

func init() {
	rootCmd.AddCommand(addrCmd())
}

func addrCmd() *cobra.Command {
	addrCmd := &cobra.Command{
		Use:     "address",
		Aliases: []string{"a", "ad", "add", "addr"},
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listAddrs)
		},
	}
	addrCmd.AddCommand(&cobra.Command{
		Use:     "list",
		Aliases: []string{"lst", "show"},
		Run: func(cmd *cobra.Command, args []string) {
			cli.runCmd(cli.listAddrs)
		},
	})
	return addrCmd
}

func (c *client) listAddrs() {
	ipcli := ip.NewWithConn(c.conn)
	links, err := c.getLinks(ipcli)
	if err != nil {
		fmt.Println("failed to get interfaces link information, err:", err)
		return
	}

	entries, err := ipcli.ListAddresses()
	if err != nil {
		fmt.Println("failed to list address entries, err:", err)
		return
	}

	ifindexs := make([]int, 0, len(links))
	for ifindex := range links {
		ifindexs = append(ifindexs, ifindex)
	}
	sort.Ints(ifindexs)

	for _, ifindex := range ifindexs {
		link := links[ifindex]
		printLinkEntry(link)
		if addrs, ok := entries[ifindex]; ok {
			for _, addr := range addrs {
				printAddrEntry(addr)
			}
		}
	}
}

func printAddrEntry(addr *ip.AddrEntry) {
	var s strings.Builder

	inet := "inet"
	if addr.Family == syscall.AF_INET6 {
		inet = "inet6"
	}
	s.WriteString(fmt.Sprintf("    %s ", inet))

	if addr.LocalAddr == nil {
		addr.LocalAddr = addr.InterfaceAddr
	}
	if addr.InterfaceAddr == nil {
		addr.InterfaceAddr = addr.LocalAddr
	}
	if addr.LocalAddr != nil {
		s.WriteString(addr.LocalAddr.String())
		if !addr.LocalAddr.Equal(addr.InterfaceAddr) {
			s.WriteString(" peer ")
			s.WriteString(addr.InterfaceAddr.String())
		}
		s.WriteString(fmt.Sprintf("/%d ", addr.PrefixLen))
	}
	if addr.BroadcastAddr != nil {
		s.WriteString(fmt.Sprintf("brd %s ", addr.BroadcastAddr.String()))
	}
	if addr.AnycastAddr != nil {
		s.WriteString(fmt.Sprintf("any %s ", addr.AnycastAddr.String()))
	}

	if addr.Scope != -1 {
		s.WriteString(fmt.Sprintf("scope %s ", addr.Scope.String()))
	}

	flags := addr.AddrFlags
	if flags == -1 {
		flags = addr.Flags
	}
	if addr.Family == syscall.AF_INET6 {
		s.WriteString(ip.AddrFlagV6(flags).String())
	} else {
		s.WriteString(flags.String())
	}

	s.WriteString(addr.Label)

	if addr.AddrInfo != nil {
		s.WriteString("\n")
		s.WriteString("       valid_lft ")
		i := addr.AddrInfo
		if i.Valid == INFINITY_LIFE_TIME {
			s.WriteString("forever")
		} else {
			s.WriteString(fmt.Sprintf("%dsec", i.Valid))
		}
		s.WriteString(" preferred_lft ")
		if i.Prefered == INFINITY_LIFE_TIME {
			s.WriteString("forever")
		} else {
			s.WriteString(fmt.Sprintf("%dsec", i.Prefered))
		}
	}

	fmt.Println(s.String())
}
