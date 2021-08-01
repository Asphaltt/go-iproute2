package bridge

import (
	"net"
	"unsafe"

	"github.com/Asphaltt/go-iproute2"
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
	State   iproute2.NudState
	Flag    iproute2.NtfFlag
	Ifindex int
	Lladdr  net.HardwareAddr
	Vlan    int
	Master  int
}

// ListFdb dumps the fdb records from the kernel.
func (c *Client) ListFdb() ([]*FdbEntry, error) {
	var ndmsg iproute2.NdMsg
	ndmsg.Family = iproute2.PF_BRIDGE

	var msg netlink.Message
	msg.Header.Type = iproute2.RTM_GETNEIGH
	msg.Header.Flags = netlink.Dump | netlink.Request
	msg.Data, _ = ndmsg.MarshalBinary()

	msgs, err := c.conn.Execute(msg)
	if err != nil {
		return nil, err
	}

	var entries []*FdbEntry
	for _, msg := range msgs {
		e, ok, err := parseFdbMsg(msg)
		if err != nil {
			return nil, err
		}
		if ok {
			entries = append(entries, e)
		}
	}
	return entries, nil
}

func DialFdbMonitor() (*netlink.Conn, error) {
	nlcfg := &netlink.Config{
		Groups: iproute2.RTNLGRP_NEIGH,
	}
	conn, err := netlink.Dial(iproute2.NETLINK_ROUTE, nlcfg)
	return conn, err
}

// MonitorFdb monitors bridge fdb entry's adding and deleting.
func (c *Client) MonitorFdb(fdbHandler func(*FdbEntry)) error {
	// join the neighbour group.
	// see: https://elixir.bootlin.com/linux/latest/source/net/bridge/br_fdb.c#L729
	//      fdb_notify function.
	_ = c.conn.JoinGroup(iproute2.RTNLGRP_NEIGH)
	defer c.conn.LeaveGroup(iproute2.RTNLGRP_NEIGH)

	for {
		msgs, err := c.conn.Receive()
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
	if msg.Header.Type != iproute2.RTM_NEWNEIGH &&
		msg.Header.Type != iproute2.RTM_DELNEIGH {
		return nil, false, nil
	}

	// the data in message is a ndmsg with following some rtattr.
	// see: https://git.kernel.org/pub/scm/network/iproute2/iproute2.git/tree/bridge/fdb.c#n137
	//      print_fdb

	data := msg.Data
	if len(data) < iproute2.SizeofNdMsg {
		return nil, false, nil
	}

	ndmsg := (*iproute2.NdMsg)(unsafe.Pointer(&data[0]))
	if ndmsg.Family != iproute2.AF_BRIDGE {
		return nil, false, nil
	}

	var entry FdbEntry
	entry.Ifindex = int(ndmsg.Ifindex)
	entry.State = iproute2.NudState(ndmsg.State)
	entry.Flag = iproute2.NtfFlag(ndmsg.Flags)
	if msg.Header.Type == iproute2.RTM_NEWNEIGH {
		entry.Action = FdbActionAdd
	} else {
		entry.Action = FdbActionDel
	}

	data = data[iproute2.SizeofNdMsg:]
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
		switch iproute2.NdAttrType(ad.Type()) {
		case iproute2.NdaLladdr:
			entry.Lladdr = net.HardwareAddr(ad.Bytes())
		case iproute2.NdaVlan:
			entry.Vlan = int(ad.Uint16())
		case iproute2.NdaMaster:
			entry.Master = int(ad.Uint32())
		}
	}
	if err := ad.Err(); err != nil {
		return nil, false, err
	}
	return &entry, true, nil
}
