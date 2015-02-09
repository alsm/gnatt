package gateway

import (
	"net"
)

type Gateway interface {
	Start()
	Port() int
	OnPacket(int, []byte, *net.UDPConn, *net.UDPAddr)
}
