package main

import (
	"fmt"
	. "github.com/alsm/gnatt/packets"
	"net/url"
)

func main() {
	cp := NewMessage(CONNECT).(*ConnectMessage)
	cp.CleanSession = true
	cp.ClientId = []byte("TestGoClient")
	cp.Duration = 30
	cp.ProtocolId = 0x01

	address, _ := url.Parse("udp://127.0.0.1:1884")
	conn, err := connectGateway(address)
	if err != nil {
		panic(err)
	}
	cp.Write(conn)
	ca, _ := ReadPacket(conn)
	fmt.Println(MessageNames[ca.MessageType()], ca.(*ConnackMessage).ReturnCode)
	reg := NewMessage(REGISTER).(*RegisterMessage)
	reg.MessageId = 0x01
	reg.TopicName = []byte("gotest")
	reg.Write(conn)
	rega, _ := ReadPacket(conn)
	fmt.Println(MessageNames[rega.MessageType()], rega.(*RegackMessage).ReturnCode, rega.(*RegackMessage).TopicId)
	tid := rega.(*RegackMessage).TopicId
	sub := NewMessage(SUBSCRIBE).(*SubscribeMessage)
	sub.TopicIdType = 0x01
	sub.MessageId = 0x02
	sub.TopicId = tid
	sub.Write(conn)
	suba, _ := ReadPacket(conn)
	fmt.Println(MessageNames[suba.MessageType()], suba.(*SubackMessage).ReturnCode, suba.(*SubackMessage).TopicId)
	pub := NewMessage(PUBLISH).(*PublishMessage)
	pub.TopicIdType = 0x01
	pub.TopicId = tid
	pub.Data = []byte("Hello Go Test")
	pub.MessageId = 0x03
	pub.Write(conn)
	pubr, _ := ReadPacket(conn)
	fmt.Println(MessageNames[pubr.MessageType()], pubr.(*PublishMessage).TopicId, pubr.(*PublishMessage).MessageId, string(pubr.(*PublishMessage).Data))
}
