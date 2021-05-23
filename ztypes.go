package iproute2

// Copying code is to avoid code for different OS platforms.

// copied from syscall
const (
	AF_BRIDGE = 0x7
	PF_BRIDGE = AF_BRIDGE

	NDA_LLADDR = 0x2
	NDA_MASTER = 0x9
)

// copied from golang.org/x/sys/unix
const (
	NETLINK_ROUTE = 0x0

	RTM_NEWNEIGH = 0x1c
	RTM_DELNEIGH = 0x1d
	RTM_GETNEIGH = 0x1e

	RTNLGRP_NEIGH = 0x3
)

const (
	SizeofNdMsg = 0xc
)

// families for netlink socket
const (
	FamilySocketMonitoring = 0x0004
)

// message types for netlink message
const (
	MsgTypeSockDiagByFamily = 0x0014
)

// SockDiagAttrType is the type for sock diag's attribute.
type SockDiagAttrType int

// attribute types for sock diag message
const (
	InetDiagNone SockDiagAttrType = iota
	InetDiagMemInfo
	InetDiagInfo
	InetDiagVegaInfo
	InetDiagCong
	InetDiagTOS
	InetDiagTclass
	InetDiagSkMemInfo
	InetDiagShutdown

	InetDiagDctcpInfo /* request as INET_DIAG_VEGASINFO */
	InetDiagProtocol  /* response attribute only */
	InetDiagSkV6Only
	InetDiagLocals
	InetDiagPeers
	InetDiagPad
	InetDiagMark    /* only with CAP_NET_ADMIN */
	InetDiagBbrInfo /* request as INET_DIAG_VEGASINFO */
	InetDiagClassID /* request as INET_DIAG_TCLASS */
	InetDiagMD5Sig
	InetDiagUlpInfo
	InetDiagSkBpfStorages
	InetDiagCgroupID
	InetDiagSockOpt
)

// SockStateType is the type for socket's state.
type SockStateType int

// states for socket
const (
	Unknown SockStateType = iota
	Established
	SynSent
	SynRecv
	FinWait1
	FinWait2
	TimeWait
	Close
	CloseWait
	LastAck
	Listen
	Closing
	_Max

	All  SockStateType = 1<<(_Max) - 1
	Conn SockStateType = All & ^((1 << Listen) | (1 << Close) | (1 << TimeWait) | (1 << SynRecv))
)

// String gets the description string of the state.
func (s SockStateType) String() string {
	switch s {
	case Established:
		return "ESTAB"
	case SynSent:
		return "SYN-SENT"
	case SynRecv:
		return "SYN-RECV"
	case FinWait1:
		return "FIN-WAIT-1"
	case FinWait2:
		return "FIN-WAIT-2"
	case TimeWait:
		return "TIME-WAIT"
	case Close:
		return "UNCONN"
	case CloseWait:
		return "CLOSE-WAIT"
	case LastAck:
		return "LAST-ACK"
	case Listen:
		return "LISTEN"
	case Closing:
		return "CLOSING"
	default:
		return "UNKNOWN"
	}
}
