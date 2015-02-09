package gateway

import (
	"fmt"
	"net"
)

func port2str(port int) string {
	return fmt.Sprintf(":%d", port)
}

func listen(g Gateway) {
	address, err := net.ResolveUDPAddr("udp", port2str(g.Port()))
	chkerr(err)
	udpconn, err := net.ListenUDP("udp", address)
	chkerr(err)
	for {
		buffer := make([]byte, 1024)
		n, remote, err := udpconn.ReadFromUDP(buffer)
		chkerr(err)
		go g.OnPacket(n, buffer, udpconn, remote)
	}
}
