package iproute2

import (
	"encoding/binary"
	"errors"
	"reflect"
	"unsafe"
)

// size of inet diag structures.
var (
	SizeofInetDiagReq = int(unsafe.Sizeof(InetDiagReq{}))
	SizeofInetDiagMsg = int(unsafe.Sizeof(InetDiagMsg{}))
)

// An InetDiagReq is a request message for sock diag netlink.
type InetDiagReq struct {
	Family   uint8
	Protocol uint8
	Ext      uint8
	Pad      uint8
	States   uint32
	InetDiagSockID
}

// An InetDiagMsg is a response message for sock diag netlink.
type InetDiagMsg struct {
	Family  uint8
	State   uint8
	Timer   uint8
	Retrans uint8

	InetDiagSockID

	Expires uint32
	RQueue  uint32
	WQueue  uint32
	UID     uint32
	Inode   uint32
}

// An InetDiagSockID contains some info about a socket.
type InetDiagSockID struct {
	Sport   uint16 // big endian
	Dport   uint16 // big endian
	Saddr   [16]byte
	Daddr   [16]byte
	Ifindex uint32
	Cookie  [2]uint32
}

// MarshalBinary marshals an inet diag request message as byte slice.
func (req *InetDiagReq) MarshalBinary() (data []byte, err error) {
	length := SizeofInetDiagReq
	var dataSlice reflect.SliceHeader
	dataSlice.Len = length
	dataSlice.Cap = length
	dataSlice.Data = uintptr(unsafe.Pointer(req))
	data = *(*[]byte)(unsafe.Pointer(&dataSlice))

	be, offset := binary.BigEndian, 8
	be.PutUint16(data[offset:], req.Sport)
	be.PutUint16(data[offset+2:], req.Dport)
	return data, nil
}

// UnmarshalBinary unmarshals an inet diag response message from byte slice.
func (msg *InetDiagMsg) UnmarshalBinary(data []byte) error {
	length := SizeofInetDiagMsg
	if len(data) < length {
		return errors.New("InetDiagMsg: not enough data to unmarshal")
	}

	dataSlice := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	newMsg := (*InetDiagMsg)(unsafe.Pointer(dataSlice.Data))
	*msg = *newMsg

	be, offset := binary.BigEndian, 4
	msg.Sport = be.Uint16(data[offset:])
	msg.Dport = be.Uint16(data[offset+2:])

	return nil
}

// NdMsg is copied from golang.org/x/sys/unix/ztypes_linux.go
type NdMsg struct {
	Family  uint8
	Pad1    uint8
	Pad2    uint16
	Ifindex int32
	State   uint16
	Flags   uint8
	Type    uint8
}

func (m *NdMsg) MarshalBinary() ([]byte, error) {
	var dataSlice reflect.SliceHeader
	dataSlice.Len = SizeofNdMsg
	dataSlice.Cap = SizeofNdMsg
	dataSlice.Data = uintptr(unsafe.Pointer(m))
	data := *(*[]byte)(unsafe.Pointer(&dataSlice))
	return data, nil
}
