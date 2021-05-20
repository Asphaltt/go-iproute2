package iproute2

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
