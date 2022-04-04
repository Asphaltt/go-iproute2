package main

import (
	"log"
	"os"
)

var config struct {
	summary bool
	process bool

	all    bool
	v4, v6 bool
	listen bool
	tcp    bool
	udp    bool
}

func showAllSocket() {
	config.udp = true
	config.tcp = true
	config.v4 = true
	config.v6 = true
	config.listen = true
}

func isMixConfig() bool {
	isMix := func(predicts ...bool) bool {
		count := 0
		for _, p := range predicts {
			if !p {
				continue
			}

			count++
			if count == 2 {
				return true
			}
		}
		return false
	}
	return isMix(
		config.tcp,
		config.udp,
	)
}

func noSock() bool {
	notToListSockets := func(predicts ...bool) bool {
		for _, p := range predicts {
			if p {
				return false
			}
		}
		return true
	}
	return notToListSockets(
		config.tcp,
		config.udp,
	)
}

func main() {
	config.v4 = true
	config.v6 = true

	for _, arg := range os.Args[1:] {
		if arg[0] != '-' {
			continue
		}

		for _, b := range arg[1:] {
			switch b {
			case 's':
				config.summary = true
			case 'a':
				config.all = true
			case '4':
				config.v4 = true
				config.v6 = false
			case '6':
				config.v4 = false
				config.v6 = true
			case 'l':
				config.listen = true
			case 't':
				config.tcp = true
			case 'u':
				config.udp = true
			case 'p':
				config.process = true
			}
		}
	}

	if config.process {
		if err := prepareProcs(); err != nil {
			log.Fatalf("failed to prepare process, err: %v", err)
		}
	}

	if config.summary {
		showSummary()
	}

	if config.all {
		showAllSocket()
	}
	if noSock() {
		if config.summary {
			return
		}

		// list all listening sockets
		config.tcp = true
		config.udp = true
	}

	c, err := dialNetlink()
	if err != nil {
		log.Fatalf("failed to create netlink socket, err: %v", err)
	}
	defer c.Close()

	c.showSocketInfoHeader()
	if config.udp {
		c.showUDP(config.listen, config.v4, config.v6)
	}
	if config.tcp {
		c.showTCP(config.listen, config.v4, config.v6)
	}
}
