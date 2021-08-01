package ip

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/Asphaltt/go-iproute2"
	"github.com/Asphaltt/go-iproute2/internal/etc"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

// LinkRxErrors is the rx error statistics of the link.
type LinkRxErrors struct {
	Length             uint64
	RingBufferOverflow uint64
	CRC                uint64
	FrameAlign         uint64
	FifoOverrun        uint64
	MissedPacket       uint64
}

// LinkTxErrors is the tx error statistics of the link.
type LinkTxErrors struct {
	Abort     uint64
	Carrier   uint64
	Fifo      uint64
	Heartbeat uint64
	Window    uint64
}

// LinkStat is the packet statistis of the link,
// including rx error statistics and tx error statistics.
type LinkStat struct {
	RxPackets   uint64
	TxPackets   uint64
	RxBytes     uint64
	TxBytes     uint64
	RxErrors    uint64
	TxErrors    uint64
	RxDropped   uint64
	TxDropped   uint64
	MulticastRx uint64
	Collisions  uint64
	LinkRxErrors
	LinkTxErrors
}

// UnmarshalBinary gets a LinkStat from a byte slice.
func (s *LinkStat) UnmarshalBinary(data []byte) error {
	sizeof := int(unsafe.Sizeof(*s))
	if len(data) < sizeof {
		return errors.New("LinkStat: not enough data to unmarshal")
	}

	newStat := (*LinkStat)(unsafe.Pointer(&data[0]))
	*s = *newStat
	return nil
}

// LinkFlags is the type of link flag.
type LinkFlags uint32

// String returns the string description of the LinkStat.
func (f LinkFlags) String() string {
	flags := []struct {
		flag uint32
		name string
	}{
		{unix.IFF_LOOPBACK, "LOOPBACK"},
		{unix.IFF_BROADCAST, "BROADCAST"},
		{unix.IFF_POINTOPOINT, "POINTOPOINT"},
		{unix.IFF_MULTICAST, "MULTICAST"},
		{unix.IFF_NOARP, "NOARP"},
		{unix.IFF_ALLMULTI, "ALLMULTI"},
		{unix.IFF_PROMISC, "PROMISC"},
		{unix.IFF_MASTER, "MASTER"},
		{unix.IFF_SLAVE, "SLAVE"},
		{unix.IFF_DEBUG, "DEBUG"},
		{unix.IFF_DYNAMIC, "DYNAMIC"},
		{unix.IFF_AUTOMEDIA, "AUTOMEDIA"},
		{unix.IFF_PORTSEL, "PORTSEL"},
		{unix.IFF_NOTRAILERS, "NOTRAILERS"},
		{unix.IFF_UP, "UP"},
		{unix.IFF_LOWER_UP, "LOWER_UP"},
		{unix.IFF_DORMANT, "DORMANT"},
		{unix.IFF_ECHO, "ECHO"},
	}

	var s []string
	if f&unix.IFF_UP != 0 && f&unix.IFF_RUNNING == 0 {
		s = append(s, "NO-CARRIER")
	}
	fl := uint32(f)
	for _, flag := range flags {
		if fl&flag.flag != 0 {
			s = append(s, flag.name)
			fl &= ^flag.flag
		}
	}
	// if fl != 0 {
	// 	s = append(s, fmt.Sprintf("%x", fl))
	// }

	return "<" + strings.Join(s, ",") + ">"
}

// LinkType is the type of link type.
type LinkType int

// String returns the string description of the LinkType.
func (t LinkType) String() string {
	// the array is copied from iproute2 C source code
	types := []struct {
		typ  int
		name string
	}{
		{unix.ARPHRD_NETROM, "netrom"},
		{unix.ARPHRD_ETHER, "ether"},
		{unix.ARPHRD_EETHER, "eether"},
		{unix.ARPHRD_AX25, "ax25"},
		{unix.ARPHRD_PRONET, "pronet"},
		{unix.ARPHRD_CHAOS, "chaos"},
		{unix.ARPHRD_IEEE802, "ieee802"},
		{unix.ARPHRD_ARCNET, "arcnet"},
		{unix.ARPHRD_APPLETLK, "atalk"},
		{unix.ARPHRD_DLCI, "dlci"},
		{unix.ARPHRD_ATM, "atm"},
		{unix.ARPHRD_METRICOM, "metricom"},
		{unix.ARPHRD_IEEE1394, "ieee1394"},
		{unix.ARPHRD_INFINIBAND, "infiniband"},
		{unix.ARPHRD_SLIP, "slip"},
		{unix.ARPHRD_CSLIP, "cslip"},
		{unix.ARPHRD_SLIP6, "slip6"},
		{unix.ARPHRD_CSLIP6, "cslip6"},
		{unix.ARPHRD_RSRVD, "rsrvd"},
		{unix.ARPHRD_ADAPT, "adapt"},
		{unix.ARPHRD_ROSE, "rose"},
		{unix.ARPHRD_X25, "x25"},
		{unix.ARPHRD_HWX25, "hwx25"},
		{unix.ARPHRD_CAN, "can"},
		{unix.ARPHRD_PPP, "ppp"},
		{unix.ARPHRD_HDLC, "hdlc"},
		{unix.ARPHRD_LAPB, "lapb"},
		{unix.ARPHRD_DDCMP, "ddcmp"},
		{unix.ARPHRD_RAWHDLC, "rawhdlc"},
		{unix.ARPHRD_TUNNEL, "ipip"},
		{unix.ARPHRD_TUNNEL6, "tunnel6"},
		{unix.ARPHRD_FRAD, "frad"},
		{unix.ARPHRD_SKIP, "skip"},
		{unix.ARPHRD_LOOPBACK, "loopback"},
		{unix.ARPHRD_LOCALTLK, "ltalk"},
		{unix.ARPHRD_FDDI, "fddi"},
		{unix.ARPHRD_BIF, "bif"},
		{unix.ARPHRD_SIT, "sit"},
		{unix.ARPHRD_IPDDP, "ip/ddp"},
		{unix.ARPHRD_IPGRE, "gre"},
		{unix.ARPHRD_PIMREG, "pimreg"},
		{unix.ARPHRD_HIPPI, "hippi"},
		{unix.ARPHRD_ASH, "ash"},
		{unix.ARPHRD_ECONET, "econet"},
		{unix.ARPHRD_IRDA, "irda"},
		{unix.ARPHRD_FCPP, "fcpp"},
		{unix.ARPHRD_FCAL, "fcal"},
		{unix.ARPHRD_FCPL, "fcpl"},
		{unix.ARPHRD_FCFABRIC, "fcfb0"},
		{unix.ARPHRD_FCFABRIC + 1, "fcfb1"},
		{unix.ARPHRD_FCFABRIC + 2, "fcfb2"},
		{unix.ARPHRD_FCFABRIC + 3, "fcfb3"},
		{unix.ARPHRD_FCFABRIC + 4, "fcfb4"},
		{unix.ARPHRD_FCFABRIC + 5, "fcfb5"},
		{unix.ARPHRD_FCFABRIC + 6, "fcfb6"},
		{unix.ARPHRD_FCFABRIC + 7, "fcfb7"},
		{unix.ARPHRD_FCFABRIC + 8, "fcfb8"},
		{unix.ARPHRD_FCFABRIC + 9, "fcfb9"},
		{unix.ARPHRD_FCFABRIC + 10, "fcfb10"},
		{unix.ARPHRD_FCFABRIC + 11, "fcfb11"},
		{unix.ARPHRD_FCFABRIC + 12, "fcfb12"},
		{unix.ARPHRD_IEEE802_TR, "tr"},
		{unix.ARPHRD_IEEE80211, "ieee802.11"},
		{unix.ARPHRD_IEEE80211_PRISM, "ieee802.11/prism"},
		{unix.ARPHRD_IEEE80211_RADIOTAP, "ieee802.11/radiotap"},
		{unix.ARPHRD_IEEE802154, "ieee802.15.4"},
		{unix.ARPHRD_IEEE802154_MONITOR, "ieee802.15.4/monitor"},
		{unix.ARPHRD_PHONET, "phonet"},
		{unix.ARPHRD_PHONET_PIPE, "phonet_pipe"},
		{unix.ARPHRD_CAIF, "caif"},
		{unix.ARPHRD_IP6GRE, "gre6"},
		{unix.ARPHRD_NETLINK, "netlink"},
		{unix.ARPHRD_6LOWPAN, "6lowpan"},

		{unix.ARPHRD_NONE, "none"},
		{unix.ARPHRD_VOID, "void"},
	}
	for _, typ := range types {
		if typ.typ == int(t) {
			return typ.name
		}
	}
	return strconv.Itoa(int(t))
}

// LinkOperState is the type of link operstate.
type LinkOperState int

// String returns the string description of the LinkOperState.
func (s LinkOperState) String() string {
	operStates := []string{
		"UNKNOWN",
		"NOTPRESENT",
		"DOWN",
		"LOWERLAYERDOWN",
		"TESTING",
		"DORMANT",
		"UP",
	}
	if int(s) >= len(operStates) {
		return fmt.Sprintf("0x%x", int(s))
	}
	return operStates[s]
}

// LinkMode is the type of link mode.
type LinkMode uint8

// String returns the string description of the LinkMode.
func (m LinkMode) String() string {
	modes := []string{"DEFAULT", "DORMANT"}
	if int(m) >= len(modes) {
		return strconv.Itoa(int(m))
	}
	return modes[m]
}

// LinkGroup is the type of link group.
type LinkGroup int

// String returns the string description of the LinkGroup.
// The group information is from */etc/iproute2/group*.
func (g LinkGroup) String() string {
	groups, _ := etc.ReadGroup()
	return groups[int(g)]
}

// A LinkEntry contains information for the link from kernel，
// like ifindex, name, link state, link type and so on.
// It should includes all information from executing command
// `ip link list`.
type LinkEntry struct {
	DeviceType       LinkType
	DeviceFlags      LinkFlags
	Ifindex          int
	Name             string
	Master           int
	Link             int
	Namespace        int
	TxQueueCount     int
	TxQueue          int
	RxQueueCount     int
	MTU              int
	MinMTU           int
	MaxMTU           int
	OperState        LinkOperState
	Mode             LinkMode
	Group            LinkGroup
	Promiscuity      int
	MaxGSOSegs       int
	MaxGSOSize       int
	Carrier          uint8
	CarrierChanges   int
	CarrierUpCount   int
	CarrierDownCount int
	QDisc            string
	ProtoDown        uint8
	Map              []byte
	Addr             []byte
	Broadcast        []byte
	Stat             []byte
	Stat64           []byte
	XDP              uint64
	AFSpec           []byte
}

// init initiates the LinkEntry to set some fields to
// value `-1`, in order to indicates that the fields
// are not included in the netlink response message.
func (e *LinkEntry) init() {
	e.Link = -1
	e.Namespace = -1
	e.Group = -1
	e.OperState = -1
}

// ListLinks gets all links information from kernel by netlink interface.
// Firstly, send a getting link request, and receive all netlink
// response messages. Secondly, parse link information from every netlink
// response messages one by one.
func (c *Client) ListLinks() ([]*LinkEntry, error) {
	var msg netlink.Message
	msg.Header.Type = unix.RTM_GETLINK
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

	entries := make([]*LinkEntry, 0, len(msgs))
	for _, msg := range msgs {
		if msg.Header.Type != unix.RTM_NEWLINK {
			continue
		}

		e, ok, err := parseLinkMsg(&msg)
		if err != nil {
			return entries, err
		}
		if ok {
			entries = append(entries, e)
		}
	}
	return entries, nil
}

// parseLinkMsg parses a link information from a netlink message.
func parseLinkMsg(msg *netlink.Message) (*LinkEntry, bool, error) {
	var ifimsg iproute2.IfInfoMsg
	if err := ifimsg.UnmarshalBinary(msg.Data); err != nil {
		return nil, false, err
	}

	var e LinkEntry
	e.init()
	e.Ifindex = int(ifimsg.Index)
	e.DeviceType = LinkType(ifimsg.Type)
	e.DeviceFlags = LinkFlags(ifimsg.Flags)

	ad, err := netlink.NewAttributeDecoder(msg.Data[iproute2.SizeofIfInfoMsg:])
	if err != nil {
		return &e, false, err
	}

	for ad.Next() {
		switch ad.Type() {
		case unix.IFLA_ADDRESS:
			e.Addr = ad.Bytes()
		case unix.IFLA_BROADCAST:
			e.Broadcast = ad.Bytes()
		case unix.IFLA_IFNAME:
			e.Name = ad.String()
		case unix.IFLA_MTU:
			e.MTU = int(ad.Uint32())
		case unix.IFLA_LINK:
			e.Link = int(ad.Uint32())
		case unix.IFLA_QDISC:
			e.QDisc = ad.String()
		case unix.IFLA_STATS:
			e.Stat = ad.Bytes()
			// if err := e.Stat.UnmarshalBinary(ad.Bytes()); err != nil {
			// 	return &e, false, err
			// }
		case unix.IFLA_MASTER:
			e.Master = int(ad.Uint32())
		case unix.IFLA_TXQLEN:
			e.TxQueue = int(ad.Uint32())
		case unix.IFLA_MAP:
			e.Map = ad.Bytes()
		case unix.IFLA_OPERSTATE:
			e.OperState = LinkOperState(ad.Bytes()[0])
		case unix.IFLA_LINKMODE:
			e.Mode = LinkMode(ad.Bytes()[0])
		case unix.IFLA_STATS64:
			e.Stat64 = ad.Bytes()
			// if err := e.Stat64.UnmarshalBinary(ad.Bytes()); err != nil {
			// 	return &e, false, err
			// }
		case unix.IFLA_AF_SPEC:
			e.AFSpec = ad.Bytes()
		case unix.IFLA_GROUP:
			e.Group = LinkGroup(ad.Uint32())
		case unix.IFLA_PROMISCUITY:
			e.Promiscuity = int(ad.Uint32())
		case unix.IFLA_NUM_TX_QUEUES:
			e.TxQueueCount = int(ad.Uint32())
		case unix.IFLA_CARRIER:
			e.Carrier = ad.Bytes()[0]
		case unix.IFLA_CARRIER_CHANGES:
			e.CarrierChanges = int(ad.Uint32())
		case unix.IFLA_LINK_NETNSID:
			e.Namespace = int(ad.Uint32())
		case unix.IFLA_PROTO_DOWN:
			e.ProtoDown = ad.Bytes()[0]
		case unix.IFLA_GSO_MAX_SEGS:
			e.MaxGSOSegs = int(ad.Uint32())
		case unix.IFLA_GSO_MAX_SIZE:
			e.MaxGSOSize = int(ad.Uint32())
		case unix.IFLA_XDP:
			e.XDP = ad.Uint64()
		case unix.IFLA_CARRIER_UP_COUNT:
			e.CarrierUpCount = int(ad.Uint32())
		case unix.IFLA_CARRIER_DOWN_COUNT:
			e.CarrierDownCount = int(ad.Uint32())
		case unix.IFLA_MIN_MTU:
			e.MinMTU = int(ad.Uint32())
		case unix.IFLA_MAX_MTU:
			e.MaxMTU = int(ad.Uint32())
		}
	}
	err = ad.Err()
	return &e, err == nil, err
}
