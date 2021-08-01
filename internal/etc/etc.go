package etc

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	routeTableConf    = "/etc/iproute2/rt_tables"
	routeProtocolConf = "/etc/iproute2/rt_protos"
	routeScopeConf    = "/etc/iproute2/rt_scopes"
	groupConf         = "/etc/iproute2/group"
)

func ReadRouteTables() (map[int]string, error) { return read(routeTableConf) }
func ReadRouteProtos() (map[int]string, error) { return read(routeProtocolConf) }
func ReadRouteScopes() (map[int]string, error) { return read(routeScopeConf) }
func ReadGroup() (map[int]string, error)       { return read(groupConf) }

func read(conf string) (map[int]string, error) {
	m := make(map[int]string)
	fd, err := os.Open(conf)
	if err != nil {
		return m, err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		items := strings.Split(line, "\t")
		if len(items) != 2 {
			continue
		}
		n, err := strconv.Atoi(items[0])
		if err == nil {
			m[n] = items[1]
		}
	}
	return m, nil
}
