package bridge

import "github.com/mdlayher/netlink"

// A Client can manipulate the bridge netlink interface.
type Client struct {
	conn *netlink.Conn
}

// New creates a Client which can issue bridger commands.
func New(conn *netlink.Conn) *Client {
	var c Client
	c.conn = conn
	return &c
}
