package ip

import (
	"github.com/Asphaltt/go-iproute2"
	"github.com/mdlayher/netlink"
)

// A Client can manipulate ip netlink interface.
type Client struct {
	conn *netlink.Conn
}

// New creates a Client which can issue ip commands.
func New() (*Client, error) {
	conn, err := netlink.Dial(iproute2.NETLINK_ROUTE, nil)
	if err != nil {
		return nil, err
	}

	return NewWithConn(conn), nil
}

// NewWithConn creates a Client which can issue ip commands with the given
// netlink connection.
func NewWithConn(conn *netlink.Conn) *Client {
	var c Client
	c.conn = conn
	return &c
}
