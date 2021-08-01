package ip

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/Asphaltt/go-iproute2"
	"github.com/Asphaltt/go-iproute2/internal/etc"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

type RouteType int

func (typ RouteType) String() string {
	switch typ {
	case unix.RTN_UNSPEC:
		return "none"
	case unix.RTN_UNICAST:
		return "unicast"
	case unix.RTN_LOCAL:
		return "local"
	case unix.RTN_BROADCAST:
		return "broadcast"
	case unix.RTN_ANYCAST:
		return "anycast"
	case unix.RTN_MULTICAST:
		return "multicast"
	case unix.RTN_BLACKHOLE:
		return "blackhole"
	case unix.RTN_UNREACHABLE:
		return "unreachable"
	case unix.RTN_PROHIBIT:
		return "prohibit"
	case unix.RTN_THROW:
		return "throw"
	case unix.RTN_NAT:
		return "nat"
	case unix.RTN_XRESOLVE:
		return "xresolve"
	default:
		return fmt.Sprintf("%d", typ)
	}
}

type RouteTable int

func (t RouteTable) String() string {
	names, _ := etc.ReadRouteTables()
	names[unix.RT_TABLE_DEFAULT] = "default"
	names[unix.RT_TABLE_MAIN] = "main"
	names[unix.RT_TABLE_LOCAL] = "local"

	if name, ok := names[int(t)]; ok {
		return name
	}
	return fmt.Sprintf("%d", t)
}

type RouteProtocol int

func (p RouteProtocol) String() string {
	protocols, _ := etc.ReadRouteProtos()
	protocols[unix.RTPROT_UNSPEC] = "unspec"
	protocols[unix.RTPROT_REDIRECT] = "redirect"
	protocols[unix.RTPROT_KERNEL] = "kernel"
	protocols[unix.RTPROT_BOOT] = "boot"
	protocols[unix.RTPROT_STATIC] = "static"
	protocols[unix.RTPROT_GATED] = "gated"
	protocols[unix.RTPROT_RA] = "ra"
	protocols[unix.RTPROT_MRT] = "mrt"
	protocols[unix.RTPROT_ZEBRA] = "zebra"
	protocols[unix.RTPROT_BIRD] = "bird"
	protocols[unix.RTPROT_BABEL] = "babel"
	protocols[unix.RTPROT_DNROUTED] = "dnrouted"
	protocols[unix.RTPROT_XORP] = "xorp"
	protocols[unix.RTPROT_NTK] = "ntk"
	protocols[unix.RTPROT_DHCP] = "dhcp"
	protocols[unix.RTPROT_KEEPALIVED] = "keepalived"
	protocols[unix.RTPROT_BGP] = "bgp"
	protocols[unix.RTPROT_ISIS] = "isis"
	protocols[unix.RTPROT_OSPF] = "ospf"
	protocols[unix.RTPROT_RIP] = "rip"
	protocols[unix.RTPROT_EIGRP] = "eigrp"

	if proto, ok := protocols[int(p)]; ok {
		return proto
	}
	return fmt.Sprintf("%d", p)
}

type RouteScope int

func (s RouteScope) String() string {
	scopes, _ := etc.ReadRouteScopes()
	scopes[unix.RT_SCOPE_UNIVERSE] = "global"
	scopes[unix.RT_SCOPE_NOWHERE] = "nowhere"
	scopes[unix.RT_SCOPE_HOST] = "host"
	scopes[unix.RT_SCOPE_LINK] = "link"
	scopes[unix.RT_SCOPE_SITE] = "site"

	if scope, ok := scopes[int(s)]; ok {
		return scope
	}
	return fmt.Sprintf("%d", s)
}

type RouteFlags int

const (
	RTM_F_OFFLOAD_FAILED = 0x20000000 /* route offload failed, this value
	 * is chosen to avoid conflicts with
	 * other flags defined in
	 * include/uapi/linux/ipv6_route.h
	 */
)

func (f RouteFlags) String() string {
	flags := []string{}
	if int(f)&unix.RTNH_F_DEAD != 0 {
		flags = append(flags, "dead")
	}
	if int(f)&unix.RTNH_F_ONLINK != 0 {
		flags = append(flags, "onlink")
	}
	if int(f)&unix.RTNH_F_PERVASIVE != 0 {
		flags = append(flags, "pervasive")
	}
	if int(f)&unix.RTNH_F_OFFLOAD != 0 {
		flags = append(flags, "offload")
	}
	// if int(f)&unix.RTNH_F_TRAP != 0 {
	// 	flags = append(flags, "trap")
	// }
	if int(f)&unix.RTM_F_NOTIFY != 0 {
		flags = append(flags, "notify")
	}
	if int(f)&unix.RTNH_F_LINKDOWN != 0 {
		flags = append(flags, "linkdown")
	}
	if int(f)&unix.RTNH_F_UNRESOLVED != 0 {
		flags = append(flags, "unresolved")
	}
	if int(f)&unix.RTM_F_OFFLOAD != 0 {
		flags = append(flags, "rt_offload")
	}
	if int(f)&unix.RTM_F_TRAP != 0 {
		flags = append(flags, "rt_trap")
	}
	if int(f)&RTM_F_OFFLOAD_FAILED != 0 {
		flags = append(flags, "rt_offload_failed")
	}
	return strings.Join(flags, " ")
}

type RoutePref int

const (
	ICMPV6_ROUTER_PREF_LOW     = 0x3
	ICMPV6_ROUTER_PREF_MEDIUM  = 0x0
	ICMPV6_ROUTER_PREF_HIGH    = 0x1
	ICMPV6_ROUTER_PREF_INVALID = 0x2
)

func (p RoutePref) String() string {
	prefs := map[int]string{
		ICMPV6_ROUTER_PREF_LOW:    "low",
		ICMPV6_ROUTER_PREF_MEDIUM: "medium",
		ICMPV6_ROUTER_PREF_HIGH:   "high",
	}
	return prefs[int(p)]
}

type RouteEntry struct {
	Family   int
	DstLen   int
	SrcLen   int
	Tos      int
	TableID  RouteTable
	Protocol RouteProtocol
	Scope    RouteScope
	Type     RouteType
	Flags    RouteFlags

	Daddr      net.IP
	Saddr      net.IP
	InIfindex  int
	OutIfindex int
	Gateway    net.IP
	Table      RouteTable
	Priority   int
	PrefSrc    net.IP
	Metric     int
	Pref       RoutePref
}

func (e *RouteEntry) init() {
	e.Table = -1
	e.Priority = -1
	e.Pref = -1
}

func (c *Client) ListRoutes() ([]*RouteEntry, error) {
	return c.listRoutes(syscall.AF_UNSPEC)
}

func (c *Client) ListRoutesV4() ([]*RouteEntry, error) {
	return c.listRoutes(syscall.AF_INET)
}

func (c *Client) ListRoutesV6() ([]*RouteEntry, error) {
	return c.listRoutes(syscall.AF_INET6)
}

func (c *Client) listRoutes(family uint8) ([]*RouteEntry, error) {
	var msg netlink.Message
	msg.Header.Type = unix.RTM_GETROUTE
	msg.Header.Flags = netlink.Dump | netlink.Request

	var rtmsg iproute2.RtMsg
	rtmsg.Family = family
	msg.Data, _ = rtmsg.MarshalBinary()

	msgs, err := c.conn.Execute(msg)
	if err != nil {
		return nil, err
	}

	entries := make([]*RouteEntry, 0, len(msgs))
	for _, msg := range msgs {
		if msg.Header.Type != unix.RTM_NEWROUTE {
			continue
		}

		e, ok, err := parseRouteMsg(&msg)
		if err != nil {
			return entries, err
		}
		if ok {
			entries = append(entries, e)
		}
	}
	return entries, nil
}

// parseRouteMsg parses a route information from a netlink message.
func parseRouteMsg(msg *netlink.Message) (*RouteEntry, bool, error) {
	var rtmsg iproute2.RtMsg
	if err := rtmsg.UnmarshalBinary(msg.Data); err != nil {
		return nil, false, err
	}

	var e RouteEntry
	e.init()
	e.Family = int(rtmsg.Family)
	e.DstLen = int(rtmsg.Dst_len)
	e.SrcLen = int(rtmsg.Src_len)
	e.Tos = int(rtmsg.Tos)
	e.TableID = RouteTable(rtmsg.Table)
	e.Protocol = RouteProtocol(rtmsg.Protocol)
	e.Scope = RouteScope(rtmsg.Scope)
	e.Type = RouteType(rtmsg.Type)
	e.Flags = RouteFlags(rtmsg.Flags)

	ad, err := netlink.NewAttributeDecoder(msg.Data[iproute2.SizeofRtMsg:])
	if err != nil {
		return &e, false, err
	}

	for ad.Next() {
		switch ad.Type() {
		case unix.RTA_DST:
			e.Daddr = net.IP(ad.Bytes())
		case unix.RTA_SRC:
			e.Saddr = net.IP(ad.Bytes())
		case unix.RTA_IIF:
			e.InIfindex = int(ad.Uint32())
		case unix.RTA_OIF:
			e.OutIfindex = int(ad.Uint32())
		case unix.RTA_GATEWAY:
			e.Gateway = net.IP(ad.Bytes())
		case unix.RTA_PRIORITY:
			e.Priority = int(ad.Uint32())
		case unix.RTA_PREFSRC:
			e.PrefSrc = net.IP(ad.Bytes())
		case unix.RTA_METRICS:
			e.Metric = int(ad.Uint32())
		case unix.RTA_TABLE:
			e.Table = RouteTable(ad.Uint32())
		case unix.RTA_PREF:
			e.Pref = RoutePref(ad.Uint8())
		}
	}
	err = ad.Err()
	return &e, err == nil, err
}
