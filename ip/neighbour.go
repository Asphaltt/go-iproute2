package ip

import (
	"net"
	"syscall"
	"unsafe"

	"github.com/Asphaltt/go-iproute2"
	"github.com/mdlayher/netlink"
)

// A NeighEntry contains information for the arp records from kernel.
type NeighEntry struct {
	Ifindex int
	Addr    net.IP
	Lladdr  net.HardwareAddr
	State   iproute2.NudState
}

// ListNeighbours dumps arp table from kernel.
// Firstly, send a getting neighbour request, and receive all netlink
// response messages. Secondly, parse neighbour information from every
// netlink response messages one by one.
func (c *Client) ListNeighbours() ([]*NeighEntry, error) {
	var ndmsg iproute2.NdMsg
	ndmsg.Family = syscall.AF_UNSPEC
	// ndmsg.State = uint16(0xFF & ^iproute2.NudNoArp)

	var msg netlink.Message
	msg.Header.Type = iproute2.RTM_GETNEIGH
	msg.Header.Flags = netlink.Dump | netlink.Request
	msg.Data, _ = ndmsg.MarshalBinary()

	msgs, err := c.conn.Execute(msg)
	if err != nil {
		return nil, err
	}

	entries := make([]*NeighEntry, 0, len(msgs))
	for _, msg := range msgs {
		if msg.Header.Type != iproute2.RTM_NEWNEIGH {
			continue
		}
		if msg.Header.Length < uint32(iproute2.SizeofNdMsg) {
			continue
		}

		ndmsg := (*iproute2.NdMsg)(unsafe.Pointer(&msg.Data[0]))
		if ndmsg.Family != syscall.AF_INET &&
			ndmsg.Family != syscall.AF_INET6 {
			continue
		}

		ad, err := netlink.NewAttributeDecoder(msg.Data[iproute2.SizeofNdMsg:])
		if err != nil {
			return entries, err
		}

		var e NeighEntry
		e.Ifindex = int(ndmsg.Ifindex)
		e.State = iproute2.NudState(ndmsg.State)

		for ad.Next() {
			switch iproute2.NdAttrType(ad.Type()) {
			case iproute2.NdaDst:
				e.Addr = net.IP(ad.Bytes())
			case iproute2.NdaLladdr:
				e.Lladdr = net.HardwareAddr(ad.Bytes())
			}
		}
		if err := ad.Err(); err != nil {
			return entries, err
		}

		entries = append(entries, &e)
	}
	return entries, nil
}
