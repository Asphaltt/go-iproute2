module github.com/Asphaltt/go-iproute2

go 1.16

require (
	github.com/mdlayher/netlink v1.4.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/spf13/cobra v1.1.3
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27
)

replace github.com/shirou/gopsutil => ./pkg/gopsutil
