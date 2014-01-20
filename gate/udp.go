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
	socket, err := net.ListenUDP("udp", address)
	chkerr(err)
	for {
		buffer := make([]byte, 1024)
		nbytes, remote, err := socket.ReadFromUDP(buffer)
		chkerr(err)
		go g.OnPacket(nbytes, buffer, remote)
	}
}
