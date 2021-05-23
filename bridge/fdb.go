package bridge

import (
	"net"
	"unsafe"

	"github.com/mdlayher/netlink"
)

// FdbActionType is the type for the fdb entry's action.
type FdbActionType int

const (
	// FdbActionAdd indicates a fdb entry was added.
	FdbActionAdd FdbActionType = iota
	// FdbActionDel indicates a fdb entry was deleted.
	FdbActionDel
)

// FdbEntry contains fdb messages for bridge fdb.
type FdbEntry struct {
	Action  FdbActionType
	Ifindex int
	Lladdr  net.HardwareAddr
	Master  int
}

// NdMsg is copied from golang.org/x/sys/unix/ztypes_linux.go
type NdMsg struct {
	Family  uint8
	Pad1    uint8
	Pad2    uint16
	Ifindex int32
	State   uint16
	Flags   uint8
	Type    uint8
}

const SizeofNdMsg = 0xc

// MonitorFdb monitors bridge fdb entry's adding and deleting.
func MonitorFdb(fdbHandler func(*FdbEntry)) error {
	nlcfg := &netlink.Config{
		Groups: RTNLGRP_NEIGH,
	}
	conn, err := netlink.Dial(NETLINK_ROUTE, nlcfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	// join the neighbour group.
	// see: https://elixir.bootlin.com/linux/latest/source/net/bridge/br_fdb.c#L729
	//      fdb_notify function.
	_ = conn.JoinGroup(RTNLGRP_NEIGH)
	defer conn.LeaveGroup(RTNLGRP_NEIGH)

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
}

func parseFdbMsg(msg netlink.Message) (*FdbEntry, bool, error) {
	// message type references
	// https://elixir.bootlin.com/linux/latest/source/net/bridge/br_fdb.c#L528
	// fdb_insert function calls fdb_notify with RTM_NEWNEIGH
	// https://elixir.bootlin.com/linux/latest/source/net/bridge/br_fdb.c#L196
	// fdb_delete function calls fdb_notify with RTM_DELNEIGH
	if msg.Header.Type != RTM_NEWNEIGH &&
		msg.Header.Type != RTM_DELNEIGH {
		return nil, false, nil
	}

	// the data in message is a ndmsg with following some rtattr.
	// see: https://git.kernel.org/pub/scm/network/iproute2/iproute2.git/tree/bridge/fdb.c#n137
	//      print_fdb

	data := msg.Data
	if len(data) < SizeofNdMsg {
		return nil, false, nil
	}

	ndmsg := (*NdMsg)(unsafe.Pointer(&data[0]))
	if ndmsg.Family != AF_BRIDGE {
		return nil, false, nil
	}

	var entry FdbEntry
	entry.Ifindex = int(ndmsg.Ifindex)
	if msg.Header.Type == RTM_NEWNEIGH {
		entry.Action = FdbActionAdd
	} else {
		entry.Action = FdbActionDel
	}

	data = data[SizeofNdMsg:]
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
		case NDA_LLADDR:
			entry.Lladdr = net.HardwareAddr(ad.Bytes())
		case NDA_MASTER:
			entry.Master = int(ad.Uint32())
		}
	}
	if err := ad.Err(); err != nil {
		return nil, false, err
	}
	return &entry, true, nil
}
