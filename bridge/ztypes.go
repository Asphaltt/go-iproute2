package bridge

// Copying code is to avoid code for different OS platforms.

// copied from syscall
const (
	AF_BRIDGE = 0x7

	NDA_LLADDR = 0x2
	NDA_MASTER = 0x9
)

// copied from golang.org/x/sys/unix
const (
	NETLINK_ROUTE = 0x0

	RTM_NEWNEIGH = 0x1c
	RTM_DELNEIGH = 0x1d

	RTNLGRP_NEIGH = 0x3
)
