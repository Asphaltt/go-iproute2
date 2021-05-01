package main

import (
	"fmt"
	"net"

	"github.com/Asphaltt/go-iproute2/bridge"
)

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
