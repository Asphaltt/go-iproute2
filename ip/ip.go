package ip

import "github.com/mdlayher/netlink"

// A Client can manipulate ip netlink interface.
type Client struct {
	conn *netlink.Conn
}

// New creates a Client which can issue ip commands.
func New(conn *netlink.Conn) *Client {
	var c Client
	c.conn = conn
	return &c
}
