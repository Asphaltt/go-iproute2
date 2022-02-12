package ss

import (
	"net"
	"syscall"

	"github.com/Asphaltt/go-iproute2"
	"github.com/mdlayher/netlink"
)

// A Client can manipulate ss netlink interface.
type Client struct {
	conn *netlink.Conn
}

// New creates a Client which can issue ss commands.
func New() (*Client, error) {
	conn, err := netlink.Dial(iproute2.FamilySocketMonitoring, nil)
	if err != nil {
		return nil, err
	}

	return NewWithConn(conn), nil
}

// NewWithConn creates a Client which can issue ss commands using an existing
// netlink connection.
func NewWithConn(conn *netlink.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

// Entry contains sockets' information, such as state, receive queue size,
// send queue size, local address, local port , peer address and peer port.
type Entry struct {
	State               iproute2.SockStateType
	RecvQ, SendQ        uint32
	LocalAddr, PeerAddr net.IP
	LocalPort, PeerPort int
	Ifindex             int
	Inode               int
	IsIPv4              bool
}

func (c *Client) listSockets(req *iproute2.InetDiagReq) ([]*Entry, error) {
	var msg netlink.Message
	msg.Header.Type = iproute2.MsgTypeSockDiagByFamily
	msg.Header.Flags = netlink.Dump | netlink.Request
	msg.Data, _ = req.MarshalBinary()

	msgs, err := c.conn.Execute(msg)
	if err != nil {
		return nil, err
	}

	entries := make([]*Entry, 0, len(msgs))
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

		var e Entry
		e.State = iproute2.SockStateType(diagMsg.State)
		e.RecvQ = diagMsg.RQueue
		e.SendQ = diagMsg.WQueue
		e.LocalPort = int(diagMsg.Sport)
		e.PeerPort = int(diagMsg.Dport)
		e.Inode = int(diagMsg.Inode)
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
