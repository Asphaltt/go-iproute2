package ss

import (
	"syscall"

	"github.com/Asphaltt/go-iproute2"
)

// ListTcp4Conns retrieves all tcp connections from kernel.
func (c *Client) ListTcp4Conns() ([]*Entry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET
	req.Protocol = syscall.IPPROTO_TCP
	req.States = uint32(iproute2.Conn)
	return c.listSockets(&req)
}

// ListTcp6Conns retrieves all tcp connections from kernel.
func (c *Client) ListTcp6Conns() ([]*Entry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET6
	req.Protocol = syscall.IPPROTO_TCP
	req.States = uint32(iproute2.Conn)
	return c.listSockets(&req)
}

// ListTcp4Listeners retreives all IPv4 tcp listeners from kernel.
func (c *Client) ListTcp4Listeners() ([]*Entry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET
	req.Protocol = syscall.IPPROTO_TCP
	req.States = uint32((1 << iproute2.Listen) | (1 << iproute2.Close))
	return c.listSockets(&req)
}

// ListTcp6Listeners retreives all IPv6 tcp listeners from kernel.
func (c *Client) ListTcp6Listeners() ([]*Entry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET6
	req.Protocol = syscall.IPPROTO_TCP
	req.States = uint32((1 << iproute2.Listen) | (1 << iproute2.Close))
	return c.listSockets(&req)
}
