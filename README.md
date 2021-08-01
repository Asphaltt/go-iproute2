# go-iproute2

Some net tools of iproute2 implement with Go.

## Completions

### ip

1. ip neigh list
2. ip link list
3. ip addr list
4. ip rourte list

### bridge

1. bridge monitor fdb
2. bridge fdb list

### ss

1. ss -s: print summary of socket stat
2. ss -l: list all udp and tcp listening sockets
3. ss -t: list all tcp connections
4. ss -lt: list all tcp listeners
5. ss -u: list all udp sockets, excluding the listening ones
6. ss -lu: list all udp listening sockets

All the above listings support IPv4 and IPv6 version.

## TO-DO

- [ ] ip
- [ ] bridge
- [x] ss (basically complete)

Maybe more tools, like dcb, devlink, nstat, tc, will be implemented in some day.
