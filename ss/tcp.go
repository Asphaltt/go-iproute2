package ss

import (
	"fmt"
	"net"
	"syscall"

	"github.com/Asphaltt/go-iproute2"
	"github.com/mdlayher/netlink"
)

// TcpEntry contains tcp sockets' information, such as state, receive queue size,
// send queue size, local address, local port , peer address and peer port.
type TcpEntry struct {
	State               iproute2.SockStateType
	RecvQ, SendQ        uint32
	LocalAddr, PeerAddr net.IP
	LocalPort, PeerPort int
	IsIPv4              bool
	Ifindex             int
}

func isZeroIPv6Addr(addr net.IP) bool {
	if addr.To16() == nil {
		return false
	}

	for i := 0; i < len(addr); i++ {
		if addr[i] != 0 {
			return false
		}
	}
	return true
}

// String formats TcpEntry into
// 'State     Recv-Q     Send-Q     Local Address:Port     Peer Address:Port'.
func (e *TcpEntry) String() string {
	var laddr string
	if e.Ifindex == 0 {
		if !e.IsIPv4 && isZeroIPv6Addr(e.LocalAddr) {
			laddr = "[::]"
		} else {
			laddr = e.LocalAddr.String()
		}
	} else {
		ifc, _ := net.InterfaceByIndex(e.Ifindex)
		if !e.IsIPv4 && isZeroIPv6Addr(e.LocalAddr) {
			laddr = fmt.Sprintf("*%%%s", ifc.Name)
		} else {
			laddr = fmt.Sprintf("%s%%%s", e.LocalAddr.String(), ifc.Name)
		}
	}
	var peerAddr string
	if !e.IsIPv4 && isZeroIPv6Addr(e.PeerAddr) {
		peerAddr = "[::]"
	} else {
		peerAddr = e.PeerAddr.String()
	}
	var peerPort string
	if e.PeerPort == 0 {
		peerPort = "*"
	} else {
		peerPort = fmt.Sprintf("%-5d", e.PeerPort)
	}
	return fmt.Sprintf(fmtEntry,
		e.State.String(),
		e.RecvQ, e.SendQ,
		laddr, e.LocalPort,
		peerAddr, peerPort,
	)
}

// ListTcpConns retrieves all tcp connections from kernel.
func (c *Client) ListTcpConns() ([]*TcpEntry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET
	req.Protocol = syscall.IPPROTO_TCP
	req.States = uint32(iproute2.Conn)
	return c.listTcpSockets(&req)
}

// ListTcp4Listeners retreives all IPv4 tcp listeners from kernel.
func (c *Client) ListTcp4Listeners() ([]*TcpEntry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET
	req.Protocol = syscall.IPPROTO_TCP
	req.States = uint32((1 << iproute2.Listen) | (1 << iproute2.Close))
	return c.listTcpSockets(&req)
}

// ListTcp6Listeners retreives all IPv6 tcp listeners from kernel.
func (c *Client) ListTcp6Listeners() ([]*TcpEntry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET6
	req.Protocol = syscall.IPPROTO_TCP
	req.States = uint32((1 << iproute2.Listen) | (1 << iproute2.Close))
	return c.listTcpSockets(&req)
}

func (c *Client) listTcpSockets(req *iproute2.InetDiagReq) ([]*TcpEntry, error) {
	var msg netlink.Message
	msg.Header.Type = iproute2.MsgTypeSockDiagByFamily
	msg.Header.Flags = netlink.Dump | netlink.Request
	msg.Data, _ = req.MarshalBinary()

	msgs, err := c.conn.Execute(msg)
	if err != nil {
		return nil, err
	}

	entries := make([]*TcpEntry, 0, len(msgs))
	for _, msg := range msgs {
		if msg.Header.Type != iproute2.MsgTypeSockDiagByFamily {
			continue
		}

		data := msg.Data
		var diagMsg iproute2.InetDiagMsg
		if err := diagMsg.UnmarshalBinary(data); err != nil {
			return entries, err
		}
		if diagMsg.Family != syscall.AF_INET &&
			diagMsg.Family != syscall.AF_INET6 {
			continue
		}

		var e TcpEntry
		e.State = iproute2.SockStateType(diagMsg.State)
		e.RecvQ = diagMsg.RQueue
		e.SendQ = diagMsg.WQueue
		e.LocalPort = int(diagMsg.Sport)
		e.PeerPort = int(diagMsg.Dport)
		e.Ifindex = int(diagMsg.Ifindex)
		if diagMsg.Family == syscall.AF_INET {
			e.IsIPv4 = true
			e.LocalAddr = net.IP(diagMsg.Saddr[:4])
			e.PeerAddr = net.IP(diagMsg.Daddr[:4])
		} else {
			e.LocalAddr = net.IP(diagMsg.Saddr[:])
			e.PeerAddr = net.IP(diagMsg.Daddr[:])
		}

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
