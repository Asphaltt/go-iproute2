package ip

import (
	"fmt"
	"net"
	"strings"

	"github.com/Asphaltt/go-iproute2"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

var AddrFlagDatas []AddrFlagData

func init() {
	AddrFlagDatas = []AddrFlagData{
		{
			Name:     "secondary",
			Mask:     unix.IFA_F_SECONDARY,
			ReadOnly: true,
			V6Only:   false,
		},
		{
			Name:     "temporary",
			Mask:     unix.IFA_F_SECONDARY,
			ReadOnly: true,
			V6Only:   false,
		},
		{
			Name:     "nodad",
			Mask:     unix.IFA_F_NODAD,
			ReadOnly: false,
			V6Only:   true,
		},
		{
			Name:     "optimistic",
			Mask:     unix.IFA_F_OPTIMISTIC,
			ReadOnly: false,
			V6Only:   true,
		},
		{
			Name:     "dadfailed",
			Mask:     unix.IFA_F_DADFAILED,
			ReadOnly: true,
			V6Only:   true,
		},
		{
			Name:     "home",
			Mask:     unix.IFA_F_HOMEADDRESS,
			ReadOnly: false,
			V6Only:   true,
		},
		{
			Name:     "deprecated",
			Mask:     unix.IFA_F_DEPRECATED,
			ReadOnly: true,
			V6Only:   true,
		},
		{
			Name:     "tentative",
			Mask:     unix.IFA_F_TENTATIVE,
			ReadOnly: true,
			V6Only:   true,
		},
		{
			Name:     "permanent",
			Mask:     unix.IFA_F_PERMANENT,
			ReadOnly: true,
			V6Only:   true,
		},
		{
			Name:     "mngtmpaddr",
			Mask:     unix.IFA_F_MANAGETEMPADDR,
			ReadOnly: false,
			V6Only:   true,
		},
		{
			Name:     "noprefixroute",
			Mask:     unix.IFA_F_NOPREFIXROUTE,
			ReadOnly: false,
			V6Only:   false,
		},
		{
			Name:     "autojoin",
			Mask:     unix.IFA_F_MCAUTOJOIN,
			ReadOnly: false,
			V6Only:   false,
		},
		{
			Name:     "stable-privacy",
			Mask:     unix.IFA_F_STABLE_PRIVACY,
			ReadOnly: true,
			V6Only:   true,
		},
	}
}

type AddrFlagData struct {
	Name     string
	Mask     int
	ReadOnly bool
	V6Only   bool
}

type AddrFlag int

func (f AddrFlag) String() string {
	var s strings.Builder
	flags := int(f)
	for _, data := range AddrFlagDatas {
		if data.Mask == unix.IFA_F_PERMANENT {
			if flags&data.Mask != 0 {
				s.WriteString("dynamic ")
			}
		} else if flags&data.Mask != 0 {
			s.WriteString(fmt.Sprintf("%s ", data.Name))
		}

		flags &= ^data.Mask
	}
	if flags != 0 {
		s.WriteString(fmt.Sprintf("flags %02x", int(f)))
	}
	return s.String()
}

type AddrFlagV6 int

func (f AddrFlagV6) String() string {
	var s strings.Builder
	flags := int(f)
	for _, data := range AddrFlagDatas {
		if data.Mask == unix.IFA_F_PERMANENT {
			if flags&data.Mask != 0 {
				s.WriteString("dynamic ")
			}
		} else if flags&data.Mask != 0 {
			if data.Mask == unix.IFA_F_SECONDARY {
				s.WriteString("temporary ")
			} else {
				s.WriteString(fmt.Sprintf("%s ", data.Name))
			}
		}

		flags &= ^data.Mask
	}
	if flags != 0 {
		s.WriteString(fmt.Sprintf("flags %02x", int(f)))
	}
	return s.String()
}

type AddrScope int

func (s AddrScope) String() string {
	scopes := map[int]string{}
	scopes[unix.RT_SCOPE_UNIVERSE] = "global"
	scopes[unix.RT_SCOPE_NOWHERE] = "nowhere"
	scopes[unix.RT_SCOPE_HOST] = "host"
	scopes[unix.RT_SCOPE_LINK] = "link"
	scopes[unix.RT_SCOPE_SITE] = "site"

	if str, ok := scopes[int(s)]; ok {
		return str
	}
	return fmt.Sprintf("%d", s)
}

type AddrEntry struct {
	Family        int
	PrefixLen     int
	Flags         AddrFlag
	Scope         AddrScope
	Ifindex       int
	Name          string
	Label         string
	InterfaceAddr net.IP
	LocalAddr     net.IP
	BroadcastAddr net.IP
	AnycastAddr   net.IP
	MulticastAddr net.IP
	AddrFlags     AddrFlag
	AddrInfo      *iproute2.IfaCacheinfo
}

func (e *AddrEntry) init() {
	e.Scope = -1
	e.AddrFlags = -1
}

// ListAddresses gets all addresses information of links from kernel
// by netlink interface.
// Firstly, send a getting address request, and receive all netlink
// response messages. Secondly, parse address information from every netlink
// response messages one by one.
func (c *Client) ListAddresses() (map[int][]*AddrEntry, error) {
	var msg netlink.Message
	msg.Header.Type = unix.RTM_GETADDR
	msg.Header.Flags = netlink.Dump | netlink.Request

	var ifimsg iproute2.IfInfoMsg
	ae := netlink.NewAttributeEncoder()
	ae.Uint32(unix.IFLA_EXT_MASK, uint32(iproute2.RTEXT_FILTER_BRVLAN))
	msg.Data, _ = ifimsg.MarshalBinary()
	data, err := ae.Encode()
	if err != nil {
		return nil, err
	}
	msg.Data = append(msg.Data, data...)

	msgs, err := c.conn.Execute(msg)
	if err != nil {
		return nil, err
	}

	entries := make(map[int][]*AddrEntry)
	for _, msg := range msgs {
		if msg.Header.Type != unix.RTM_NEWADDR {
			continue
		}

		e, ok, err := parseAddrMsg(&msg)
		if err != nil {
			return entries, err
		}
		if ok {
			entries[e.Ifindex] = append(entries[e.Ifindex], e)
		}
	}
	return entries, nil
}

// parseAddrMsg parses a link address information from a netlink message.
func parseAddrMsg(msg *netlink.Message) (*AddrEntry, bool, error) {
	var ifamsg iproute2.IfAddrMsg
	if err := ifamsg.UnmarshalBinary(msg.Data); err != nil {
		return nil, false, err
	}

	var e AddrEntry
	e.init()
	e.Family = int(ifamsg.Family)
	e.PrefixLen = int(ifamsg.Prefixlen)
	e.Flags = AddrFlag(ifamsg.Flags)
	e.Scope = AddrScope(ifamsg.Scope)
	e.Ifindex = int(ifamsg.Index)

	ad, err := netlink.NewAttributeDecoder(msg.Data[iproute2.SizeofIfAddrMsg:])
	if err != nil {
		return &e, false, err
	}

	for ad.Next() {
		switch ad.Type() {
		case unix.IFA_ADDRESS:
			e.InterfaceAddr = net.IP(ad.Bytes())
		case unix.IFA_LOCAL:
			e.LocalAddr = net.IP(ad.Bytes())
		case unix.IFA_LABEL:
			e.Label = ad.String()
		case unix.IFA_BROADCAST:
			e.BroadcastAddr = net.IP(ad.Bytes())
		case unix.IFA_ANYCAST:
			e.AnycastAddr = net.IP(ad.Bytes())
		case unix.IFA_CACHEINFO:
			e.AddrInfo = new(iproute2.IfaCacheinfo)
			_ = e.AddrInfo.UnmarshalBinary(ad.Bytes())
		case unix.IFA_MULTICAST:
			e.MulticastAddr = net.IP(ad.Bytes())
		case unix.IFA_FLAGS:
			e.AddrFlags = AddrFlag(ad.Uint32())
		}
	}
	err = ad.Err()
	return &e, err == nil, err
}
