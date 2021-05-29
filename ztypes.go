package iproute2

import "math/bits"

// Copying code is to avoid code for different OS platforms.

// copied from syscall
const (
	AF_BRIDGE = 0x7
	PF_BRIDGE = AF_BRIDGE
)

// copied from golang.org/x/sys/unix
const (
	NETLINK_ROUTE = 0x0

	RTM_NEWNEIGH = 0x1c
	RTM_DELNEIGH = 0x1d
	RTM_GETNEIGH = 0x1e

	RTNLGRP_NEIGH = 0x3
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

type NdAttrType uint16

// types for error message attribute
const (
	NdaUnspec NdAttrType = iota
	NdaDst
	NdaLladdr
	NdaCacheInfo
	NdaProbes
	NdaVlan
	NdaPort
	NdaVNI
	NdaIfindex
	NdaMaster
	NdaLinkNetNSID
	NdaSrcVNI
	NdaProtocol
	NdaNhID
	NdaFdbExtAttrs
)

type NtfFlag uint8

const (
	NtfUse NtfFlag = 1 << iota
	NtfSelf
	NtfMaster
	NtfProxy
	NtfExtLearned
	NtfOffloaded
	NtfSticky
	NtfRouter
)

func (f NtfFlag) String() string {
	if f == 0 {
		return ""
	}
	flags := [...]string{
		"use",
		"self",
		"master",
		"proxy",
		"extern_learn",
		"offload",
		"sticky",
		"router",
	}
	index := bits.TrailingZeros8(uint8(f))
	return flags[index]
}

type NudState uint16

const (
	NudIncomplete NudState = 1 << iota
	NudReachable
	NudStale
	NudDelay
	NudProbe
	NudFailed
	NudNoArp
	NudPermanent
	NudNone NudState = 0
)

func (s NudState) String() string {
	if s == NudNone {
		return "none"
	}
	states := [...]string{
		"incomplete",
		"reachable",
		"stale",
		"delay",
		"probe",
		"failed",
		"noarp",
		"permanent",
	}
	index := bits.TrailingZeros16(uint16(s))
	return states[index]
}

// copied from src/cmd/vendor/golang.org/x/sys/unix/ztypes_linux.go
// TODO(Asphaltt): use ztypes_linux.go instead
const (
	IFLA_EXT_MASK = 0x1d
)

/* New extended info filters for IFLA_EXT_MASK */
type RtTextFilterType int

const (
	RTEXT_FILTER_VF RtTextFilterType = 1 << iota
	RTEXT_FILTER_BRVLAN
	RTEXT_FILTER_BRVLAN_COMPRESSED
	RTEXT_FILTER_SKIP_STATS
	RTEXT_FILTER_MRP
	RTEXT_FILTER_CFM_CONFIG
	RTEXT_FILTER_CFM_STATUS
)
