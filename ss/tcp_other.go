//+build !linux

package ss

import "github.com/mdlayher/netlink"

// ListTcpConns retrieves all tcp connections from the netlink `conn`.
func ListTcpConns(conn *netlink.Conn) ([]*TcpEntry, error) {
	return nil, ErrNotImplemented
}
