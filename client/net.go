package main

import (
	"errors"
	//. "github.com/alsm/gnatt/packets"
	"net"
	"net/url"
)

func connectGateway(uri *url.URL) (conn net.Conn, err error) {
	if uri.Scheme != "udp" {
		err = errors.New("Unsupported transport")
		return
	} else {
		conn, err = net.Dial("udp", uri.Host)
	}
	return
}
