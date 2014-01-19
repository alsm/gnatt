package gateway

import (
	"fmt"
	"github.com/alsm/gnatt/common/utils"
	"net"
)

type udpListener struct {
	port int
}

func port2str(port int) string {
	return fmt.Sprintf(":%d", port)
}

func (udp *udpListener) listen() {
	address, err := net.ResolveUDPAddr("udp", port2str(udp.port))
	chkerr(err)
	socket, err := net.ListenUDP("udp", address)
	chkerr(err)
	for {
		buffer := make([]byte, 1024)
		nbytes, remote, err := socket.ReadFromUDP(buffer)
		chkerr(err)
		go handler(nbytes, buffer, remote)
	}
}

func handler(nbytes int, packet []byte, remote *net.UDPAddr) {
	fmt.Printf("received %d bytes from remote %v\n", nbytes, remote)
	fmt.Printf(" %s\n", utils.Bytes2str(packet[0:nbytes]))
}
