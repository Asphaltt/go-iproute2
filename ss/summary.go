package ss

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	procSockStat  = "/proc/net/sockstat"
	procSockStat6 = "/proc/net/sockstat6"
	procSnmp      = "/proc/net/snmp"
)

// Summary is the kernel sockets' statistics.
type Summary struct {
	Sockets    int
	TcpMem     int
	TcpTotal   int
	TcpOrphans int
	TcpTws     int
	Tcp4Hashed int
	Udp4       int
	Raw4       int
	Frag4      int
	Frag4Mem   int
	Tcp6Hashed int
	Udp6       int
	Raw6       int
	Frag6      int
	Frag6Mem   int

	TcpEstablished int
}

func (_ *Summary) parseNumber(s string, field *int) {
	n, err := strconv.Atoi(s)
	if err == nil {
		*field = n
	}
}

func (s *Summary) parseLine(line string, fields ...*int) {
	elems := strings.Split(line, " ")
	if len(elems) != 2*len(fields)+1 {
		return
	}
	for i := 1; i <= len(fields); i++ {
		s.parseNumber(elems[2*i], fields[i-1])
	}
}

// GetSockStat reads IPv4 socket statistics data from
// */proc/net/sockstat*.
func (s *Summary) GetSockStat() error {
	return s.getSockStat(procSockStat)
}

// GetSockStat6 reads IPv6 socket statistics data from
// */proc/net/sockstat6*.
func (s *Summary) GetSockStat6() error {
	return s.getSockStat(procSockStat6)
}

func (s *Summary) getSockStat(stat string) error {
	fd, err := os.Open(stat)
	if err != nil {
		return err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "sockets:") {
			s.parseLine(line, &s.Sockets)
		} else if strings.HasPrefix(line, "TCP:") {
			s.parseLine(line, &s.Tcp4Hashed,
				&s.TcpOrphans, &s.TcpTws,
				&s.TcpTotal, &s.TcpMem)
		} else if strings.HasPrefix(line, "TCP6:") {
			s.parseLine(line, &s.Tcp6Hashed)
		} else if strings.HasPrefix(line, "UDP:") {
			s.parseLine(line, &s.Udp4)
		} else if strings.HasPrefix(line, "UDP6:") {
			s.parseLine(line, &s.Udp6)
		} else if strings.HasPrefix(line, "RAW:") {
			s.parseLine(line, &s.Raw4)
		} else if strings.HasPrefix(line, "RAW6:") {
			s.parseLine(line, &s.Raw6)
		} else if strings.HasPrefix(line, "FRAG:") {
			s.parseLine(line, &s.Frag4, &s.Frag4Mem)
		} else if strings.HasPrefix(line, "FRAG6:") {
			s.parseLine(line, &s.Frag6, &s.Frag6Mem)
		}
	}
	return scanner.Err()
}

// GetTcpEstablished gets the number of established tcp connections
// from */proc/net/snmp*.
func (s *Summary) GetTcpEstablished() error {
	fd, err := os.Open(procSnmp)
	if err != nil {
		return err
	}
	defer fd.Close()

	const fieldName = "CurrEstab"
	var fieldIndex = 0

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "Tcp:") {
			continue
		}

		// Tcp: RtoAlgorithm RtoMin RtoMax MaxConn ActiveOpens PassiveOpens AttemptFails EstabResets CurrEstab InSegs OutSegs RetransSegs InErrs OutRsts InCsumErrors
		// Tcp: 1 200 120000 -1 139 112105 7639 939 3 1624441 1927238 536033 0 12449 0
		elems := strings.Split(line, " ")
		if fieldIndex == 0 {
			for fieldIndex = 1; fieldIndex < len(elems) && elems[fieldIndex] != fieldName; fieldIndex++ {
			}
		} else if fieldIndex < len(elems) {
			s.parseNumber(elems[fieldIndex], &s.TcpEstablished)
			break
		}
	}
	return scanner.Err()
}
