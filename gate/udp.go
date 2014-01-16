package gateway

import (
	"fmt"
	"net"
)

type udpListener struct {
	port int
}

func port2str(port int) string {
	return fmt.Sprintf(":%d", port)
}

func (udp *udpListener) listen() {
	buffer := make([]byte, 1024)
	address, err := net.ResolveUDPAddr("udp", port2str(udp.port))
	chkerr(err)
	socket, err := net.ListenUDP("udp", address)
	chkerr(err)
	for {
		nbytes, remote, err := socket.ReadFromUDP(buffer)
		chkerr(err)
		fmt.Printf("received %d bytes from remote %v\n", nbytes, remote)
		fmt.Printf("  %s\n", string(buffer))
	}
}
