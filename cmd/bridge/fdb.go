package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/Asphaltt/go-iproute2/bridge"
)

func listFdb() {
	entries, err := bridge.ListFdb()
	if err != nil {
		fmt.Println("failed to list fdb entries, err:", err)
		return
	}

	for _, e := range entries {
		printListFdb(e)
	}
}

func printListFdb(e *bridge.FdbEntry) {
	var devName string
	if ifc, err := net.InterfaceByIndex(e.Ifindex); err == nil {
		devName = ifc.Name
	}

	// TODO: show vlan
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s dev %s", e.Lladdr, devName))
	if e.Flag != 0 {
		b.WriteString(fmt.Sprintf(" %s", e.Flag))
	}
	b.WriteString(fmt.Sprintf(" %s", e.State))
	fmt.Println(b.String())
}

func monitorFdb() {
	err := bridge.MonitorFdb(printFdbEntry)
	if err != nil {
		fmt.Println("failed to bridge monitor fdb, err:", err)
	}
}

func printFdbEntry(entry *bridge.FdbEntry) {
	var action string
	switch entry.Action {
	case bridge.FdbActionAdd:
		action = "Added"
	case bridge.FdbActionDel:
		action = "Deleted"
	default:
		action = "Unkowned"
	}
	devInfo, _ := net.InterfaceByIndex(entry.Ifindex)
	masterInfo, _ := net.InterfaceByIndex(entry.Master)
	fmt.Printf("%s %s dev %s master %s\n", action, entry.Lladdr, devInfo.Name, masterInfo.Name)
}
