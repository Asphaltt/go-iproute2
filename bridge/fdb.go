package bridge

import (
	"net"
)

// FdbActionType is the type for the fdb entry's action.
type FdbActionType int

const (
	// FdbActionAdd indicates a fdb entry was added.
	FdbActionAdd FdbActionType = iota
	// FdbActionDel indicates a fdb entry was deleted.
	FdbActionDel
)

// FdbEntry contains fdb messages for bridge fdb.
type FdbEntry struct {
	Action  FdbActionType
	Ifindex int
	Lladdr  net.HardwareAddr
	Master  int
}
