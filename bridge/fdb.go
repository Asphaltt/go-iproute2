package bridge

import (
	"encoding/binary"
	"net"

	goendian "github.com/virtao/GoEndian"
)

var (
	endian binary.ByteOrder
)

func init() {
	if goendian.IsBigEndian() {
		endian = binary.BigEndian
	} else {
		endian = binary.LittleEndian
	}
}

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
