package main

import (
	"fmt"
	"net"
)

type udpListener struct {
}

func (gate *udpListener) listen() {
	buffer := make([]byte, 1024)
	address, err := net.ResolveUDPAddr("udp", ":16111")
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
