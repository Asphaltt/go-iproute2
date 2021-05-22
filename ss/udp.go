package ss

import (
	"syscall"

	"github.com/Asphaltt/go-iproute2"
)

// ListUdp4Conns retrieves all Udp sockets from kernel.
func (c *Client) ListUdp4Sockets() ([]*Entry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET
	req.Protocol = syscall.IPPROTO_UDP
	req.States = uint32(1 << iproute2.Established)
	return c.listSockets(&req)
}

// ListUdp6Conns retrieves all Udp sockets from kernel.
func (c *Client) ListUdp6Sockets() ([]*Entry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET6
	req.Protocol = syscall.IPPROTO_UDP
	req.States = uint32(1 << iproute2.Established)
	return c.listSockets(&req)
}

// ListUdp4Listeners retreives all IPv4 Udp listeners from kernel.
func (c *Client) ListUdp4Listeners() ([]*Entry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET
	req.Protocol = syscall.IPPROTO_UDP
	req.States = uint32((1 << iproute2.Listen) | (1 << iproute2.Close))
	return c.listSockets(&req)
}

// ListUdp6Listeners retreives all IPv6 Udp listeners from kernel.
func (c *Client) ListUdp6Listeners() ([]*Entry, error) {
	var req iproute2.InetDiagReq
	req.Family = syscall.AF_INET6
	req.Protocol = syscall.IPPROTO_UDP
	req.States = uint32((1 << iproute2.Listen) | (1 << iproute2.Close))
	return c.listSockets(&req)
}
