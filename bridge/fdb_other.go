//+build !linux

package bridge

// MonitorFdb monitors bridge fdb entry's adding and deleting.
func MonitorFdb(fdbHandler func(*FdbEntry)) error {
	return ErrNotImplemented
}
