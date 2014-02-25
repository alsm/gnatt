package gateway

import (
	"fmt"
	"os"

	. "github.com/alsm/gnatt/common/protocol"
	"github.com/alsm/gnatt/common/utils"
)

type TransGate struct {
	stopsig    chan os.Signal
	port       int
	mqttbroker string
}

func NewTransGate(gc *GatewayConfig, stopsig chan os.Signal) *TransGate {
	var tg TransGate
	tg.port = gc.port
	tg.mqttbroker = gc.mqttbroker
	return &tg
}

func (tg *TransGate) Port() int {
	return tg.port
}

func (tg *TransGate) Start() {
	go tg.awaitStop()
	fmt.Println("Transparent Gateway is starting")
	fmt.Println("Transparent Gataway is started")
	listen(tg)
}

func (tg *TransGate) awaitStop() {
	<-tg.stopsig
	fmt.Println("Transparent Gateway is stopping")
	fmt.Println("Transparent Gateway is stopped")
	os.Exit(0)
}

func (tg *TransGate) OnPacket(nbytes int, buffer []byte, con uConn, addr uAddr) {
	fmt.Println("TG OnPacket!")
	fmt.Printf("bytes: %s\n", utils.Bytes2str(buffer[0:nbytes]))
	m := Unpack(buffer[0:nbytes])
	fmt.Printf("m.MsgType(): %s\n", m.MsgType())
}
