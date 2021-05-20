//+build linux

package bridge

import (
	"net"
	"syscall"
	"unsafe"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

// MonitorFdb monitors bridge fdb entry's adding and deleting.
func MonitorFdb(fdbHandler func(*FdbEntry)) error {
	nlcfg := &netlink.Config{
		Groups: syscall.RTNLGRP_NEIGH,
	}
	conn, err := netlink.Dial(syscall.NETLINK_ROUTE, nlcfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	// join the neighbour group.
	// see: https://elixir.bootlin.com/linux/latest/source/net/bridge/br_fdb.c#L729
	//      fdb_notify function.
	_ = conn.JoinGroup(syscall.RTNLGRP_NEIGH)
	defer conn.LeaveGroup(syscall.RTNLGRP_NEIGH)

	for {
		msgs, err := conn.Receive()
		if err != nil {
			return err
		}

		for _, msg := range msgs {
			entry, ok, err := parseFdbMsg(msg)
			if err != nil {
				return err
			}
			if ok {
				fdbHandler(entry)
			}
		}
	}
	return nil
}

func parseFdbMsg(msg netlink.Message) (*FdbEntry, bool, error) {
	// message type references
	// https://elixir.bootlin.com/linux/latest/source/net/bridge/br_fdb.c#L528
	// fdb_insert function calls fdb_notify with RTM_NEWNEIGH
	// https://elixir.bootlin.com/linux/latest/source/net/bridge/br_fdb.c#L196
	// fdb_delete function calls fdb_notify with RTM_DELNEIGH
	if msg.Header.Type != syscall.RTM_NEWNEIGH &&
		msg.Header.Type != syscall.RTM_DELNEIGH {
		return nil, false, nil
	}

	// the data in message is a ndmsg with following some rtattr.
	// see: https://git.kernel.org/pub/scm/network/iproute2/iproute2.git/tree/bridge/fdb.c#n137
	//      print_fdb

	data := msg.Data
	if len(data) < unix.SizeofNdMsg {
		return nil, false, nil
	}

	ndmsg := (*unix.NdMsg)(unsafe.Pointer(&data[0]))
	if ndmsg.Family != syscall.AF_BRIDGE {
		return nil, false, nil
	}

	var entry FdbEntry
	entry.Ifindex = int(ndmsg.Ifindex)
	if msg.Header.Type == syscall.RTM_NEWNEIGH {
		entry.Action = FdbActionAdd
	} else {
		entry.Action = FdbActionDel
	}

	data = data[unix.SizeofNdMsg:]
	// netlink.Attribute describes all types of attribute in netlink message,
	// including struct rtattr.
	// the following attributes are rtattr.
	// referece: https://github.com/golang/go/blob/master/src/syscall/netlink_linux.go#L148
	//           ParseNetlinkRouteAttr
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return nil, false, err
	}
	for ad.Next() {
		switch ad.Type() {
		case unix.NDA_LLADDR:
			entry.Lladdr = net.HardwareAddr(ad.Bytes())
		case unix.NDA_MASTER:
			entry.Master = int(ad.Uint32())
		}
	}
	if err := ad.Err(); err != nil {
		return nil, false, err
	}
	return &entry, true, nil
}
