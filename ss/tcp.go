package ss

import (
	"fmt"
	"net"

	"github.com/Asphaltt/go-iproute2"
)

// TcpEntry contains tcp sockets' information, such as state, receive queue size,
// send queue size, local address, local port , peer address and peer port.
type TcpEntry struct {
	State               iproute2.SockStateType
	RecvQ, SendQ        uint32
	LocalAddr, PeerAddr net.IP
	LocalPort, PeerPort int
	Ifindex             int
}

// String formats TcpEntry into
// 'State     Recv-Q     Send-Q     Local Address:Port     Peer Address:Port'.
func (e *TcpEntry) String() string {
	var laddr string
	if e.Ifindex == 0 {
		laddr = e.LocalAddr.String()
	} else {
		ifc, _ := net.InterfaceByIndex(e.Ifindex)
		laddr = fmt.Sprintf("%s%%%s", e.LocalAddr.String(), ifc.Name)
	}
	return fmt.Sprintf(fmtEntry,
		e.State.String(),
		e.RecvQ, e.SendQ,
		laddr, e.LocalPort,
		e.PeerAddr, e.PeerPort,
	)
}
