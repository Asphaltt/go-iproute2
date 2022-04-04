package main

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/process"
)

type processInfo struct {
	name string
	pid  int32
	fd   uint32
}

func (p *processInfo) String() string {
	return fmt.Sprintf(`("%s",pid=%d,fd=%d)`, p.name, p.pid, p.fd)
}

type processesInfo []*processInfo

func (ps processesInfo) String() string {
	s := make([]string, 0, len(ps))
	for _, p := range ps {
		s = append(s, p.String())
	}
	return strings.Join(s, ",")
}

var procs = make(map[int]processesInfo)

func prepareProcs() error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}

	for _, p := range processes {
		conns, err := p.Connections()
		if err != nil {
			return err
		}

		pname, err := p.Name()
		if err != nil {
			return err
		}
		for _, c := range conns {
			procs[c.Inode] = append(procs[c.Inode], &processInfo{
				name: pname,
				pid:  p.Pid,
				fd:   c.Fd,
			})
		}
	}

	return nil
}
