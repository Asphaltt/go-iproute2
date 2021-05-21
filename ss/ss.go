package ss

import "github.com/mdlayher/netlink"

// "State     Recv-Q     Send-Q     Local Address:Port     Peer Address:Port"
const fmtEntry = "%-10s     %-6d     %-6d    %24s:%-5d     %24s:%s"

type Client struct {
	conn *netlink.Conn
}

func NewClient(conn *netlink.Conn) *Client {
	return &Client{
		conn: conn,
	}
}
