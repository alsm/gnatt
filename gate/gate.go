package gateway

import (
	"net"
)

type gwStatus byte

const (
	gw_stopped gwStatus = iota
	gw_starting
	gw_running
	gw_stopping
)

type Gateway interface {
	Start()
	Port() int
	OnPacket(int, []byte, *net.UDPAddr)
}
