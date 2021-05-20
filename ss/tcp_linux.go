//+build linux

package ss

import (
	"net"
	"syscall"

	"github.com/Asphaltt/go-iproute2"
	// "github.com/davecgh/go-spew/spew"
	"github.com/mdlayher/netlink"
)

// ListTcpConns retrieves all tcp connections from the netlink `conn`.
func ListTcpConns(conn *netlink.Conn) ([]*TcpEntry, error) {
	var msg netlink.Message
	msg.Header.Type = iproute2.MsgTypeSockDiagByFamily
	msg.Header.Flags = netlink.Root | netlink.Match | netlink.Request

	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET
	req.Protocol = syscall.IPPROTO_TCP
	msg.Data, _ = req.MarshalBinary()
	req.States = uint32(iproute2.Conn)

	// spew.Dump(msg.Data)

	msgs, err := conn.Execute(msg)
	if err != nil {
		return nil, err
	}

	entries := make([]*TcpEntry, 0, len(msgs))
	for _, msg := range msgs {
		if msg.Header.Type != iproute2.MsgTypeSockDiagByFamily {
			continue
		}

		// spew.Dump(msg.Data)

		data := msg.Data
		var diagMsg iproute2.InetDiagMsg
		if err := diagMsg.UnmarshalBinary(data); err != nil {
			return entries, err
		}
		if diagMsg.Family != syscall.AF_INET {
			continue
		}

		var e TcpEntry
		e.State = iproute2.SockStateType(diagMsg.State)
		e.RecvQ = diagMsg.RQueue
		e.SendQ = diagMsg.WQueue
		e.LocalAddr = net.IP(diagMsg.Saddr[:4])
		e.LocalPort = int(diagMsg.Sport)
		e.PeerAddr = net.IP(diagMsg.Daddr[:4])
		e.PeerPort = int(diagMsg.Dport)
		e.Ifindex = int(diagMsg.Ifindex)

		data = data[iproute2.SizeofInetDiagMsg:]
		ad, err := netlink.NewAttributeDecoder(data)
		if err != nil {
			return entries, err
		}
		for ad.Next() {
			switch ad.Type() {
			}
		}
		if err := ad.Err(); err != nil {
			return entries, err
		}

		entries = append(entries, &e)
	}
	return entries, nil
}
