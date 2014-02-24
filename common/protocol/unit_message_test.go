package gnatt

import (
	//"fmt"
	"testing"
)

func Test_MessageTypeValues(t *testing.T) {
	if ADVERTISE != 0x00 {
		t.Errorf("ADVERTISE is %x instead of %x", ADVERTISE, 0x00)
	}
	if SEARCHGW != 0x01 {
		t.Errorf("SEARCHGW is %x instead of %x", SEARCHGW, 0x01)
	}
	if GWINFO != 0x02 {
		t.Errorf("GWINFO is %x instead of %x", GWINFO, 0x02)
	}
	if CONNECT != 0x04 {
		t.Errorf("CONNECT is %x instead of %x", CONNECT, 0x04)
	}
	if CONNACK != 0x05 {
		t.Errorf("CONNACK is %x instead of %x", CONNACK, 0x05)
	}
	if WILLTOPICREQ != 0x06 {
		t.Errorf("WILLTOPICREQ is %x instead of %x", WILLTOPICREQ, 0x06)
	}
	if WILLTOPIC != 0x07 {
		t.Errorf("WILLTOPIC is %x instead of %x", WILLTOPIC, 0x07)
	}
	if WILLMSGREQ != 0x08 {
		t.Errorf("WILLMSGREQ is %x instead of %x", WILLMSGREQ, 0x08)
	}
	if WILLMSG != 0x09 {
		t.Errorf("WILLMSG is %x instead of %x", WILLMSG, 0x09)
	}
	if REGISTER != 0x0A {
		t.Errorf("REGISTER is %x instead of %x", REGISTER, 0x0A)
	}
	if REGACK != 0x0B {
		t.Errorf("REGACK is %x instead of %x", REGACK, 0x0B)
	}
	if PUBLISH != 0x0C {
		t.Errorf("PUBLISH is %x instead of %x", PUBLISH, 0x0C)
	}
	if PUBACK != 0x0D {
		t.Errorf("PUBACK is %x instead of %x", PUBACK, 0x0D)
	}
	if PUBCOMP != 0x0E {
		t.Errorf("PUBCOMP is %x instead of %x", PUBCOMP, 0x0E)
	}
	if PUBREC != 0x0F {
		t.Errorf("PUBREC is %x instead of %x", PUBREC, 0x0F)
	}
	if PUBREL != 0x10 {
		t.Errorf("PUBREL is %x instead of %x", PUBREL, 0x10)
	}
	if SUBSCRIBE != 0x12 {
		t.Errorf("SUBSCRIBE is %x instead of %x", SUBSCRIBE, 0x12)
	}
	if SUBACK != 0x13 {
		t.Errorf("SUBACK is %x instead of %x", SUBACK, 0x13)
	}
	if UNSUBSCRIBE != 0x14 {
		t.Errorf("UNSUBSCRIBE is %x instead of %x", UNSUBSCRIBE, 0x14)
	}
	if UNSUBACK != 0x15 {
		t.Errorf("UNSUBACK is %x instead of %x", UNSUBACK, 0x15)
	}
	if PINGREQ != 0x16 {
		t.Errorf("PINGREQ is %x instead of %x", PINGREQ, 0x16)
	}
	if PINGRESP != 0x17 {
		t.Errorf("PINGRESP is %x intead of %x", PINGRESP, 0x17)
	}
	if DISCONNECT != 0x18 {
		t.Errorf("DISCONNECT is %x intead of %x", DISCONNECT, 0x18)
	}
	if WILLTOPICUPD != 0x1A {
		t.Errorf("WILLTOPICUPD is %x intead of %x", WILLTOPICUPD, 0x1A)
	}
	if WILLTOPICRESP != 0x1B {
		t.Errorf("WILLTOPICRESPis %x intead of %x", WILLTOPICRESP, 0x1B)
	}
	if WILLMSGUPD != 0x1C {
		t.Errorf("WILLMSGUPD is %x intead of %x", WILLMSGUPD, 0x1C)
	}
	if WILLMSGRESP != 0x1D {
		t.Errorf("WILLMSGRESP is %x intead of %x", WILLMSGRESP, 0x1D)
	}
}

func Test_Length(t *testing.T) {
	var h Header
	m := map[int]int{
		0:     0x0000,
		1:     0x0001,
		2:     0x0002,
		3:     0x0003,
		8:     0x0008,
		15:    0x000F,
		16:    0x0010,
		127:   0x007F,
		256:   0x0100,
		16384: 0x4000,
		65535: 0xFFFF,
	}

	for in, exp := range m {
		h.SetLength(in)
		if h.Length() != exp {
			t.Errorf("encodeLength failed, input %d expected 0x%X, got 0x%X", in, exp, h.Length())
		}
	}
}

func Test_NewAdvertiseMessage(t *testing.T) {
	am := NewMessage(ADVERTISE)
	switch am.(type) {
	case *AdvertiseMessage:

	default:
		t.Errorf("msg should be *advertiseMessage")
	}
}

func Test_NewSearchGwMessage(t *testing.T) {
	sm := NewMessage(SEARCHGW)
	switch sm.(type) {
	case *SearchgwMessage:

	default:
		t.Errorf("msg should be *searchgwMessage")
	}
}

func Test_NewGwInfoMessage(t *testing.T) {
	gm := NewMessage(GWINFO)
	switch gm.(type) {
	case *GwInfoMessage:

	default:
		t.Errorf("msg should be *gwInfoMessage")
	}
}

func Test_NewConnectMessage(t *testing.T) {
	cm := NewMessage(CONNECT)
	switch cm.(type) {
	case *ConnectMessage:

	default:
		t.Errorf("msg should be *connectMessage")
	}
}

func Test_NewConnackMessage(t *testing.T) {
	cm := NewMessage(CONNACK)
	switch cm.(type) {
	case *ConnackMessage:

	default:
		t.Errorf("msg should be *connackMessage")
	}
}

func Test_NewWillTopicReqMessage(t *testing.T) {
	wm := NewMessage(WILLTOPICREQ)
	switch wm.(type) {
	case *WillTopicReqMessage:

	default:
		t.Errorf("msg should be *willTopicReqMessage")
	}
}

func Test_NewWillTopicMessage(t *testing.T) {
	m := NewMessage(WILLTOPIC)
	switch wtm := m.(type) {
	case *WillTopicMessage:
		if wtm.QoS() != QoS_Zero {
			t.Errorf("default QoS should be QoS_Zero, got %d", wtm.QoS())
		}

	default:
		t.Errorf("wtm should be *willTopicMessage")
	}
}

func Test_NewWillMsgReqMessage(t *testing.T) {
	wm := NewMessage(WILLMSGREQ)
	switch wm.(type) {
	case *WillMsgReqMessage:

	default:
		t.Errorf("msg should be *willMsgReqMessage")
	}
}

func Test_NewWillMsgMessage(t *testing.T) {
	wm := NewMessage(WILLMSG)
	switch wm.(type) {
	case *WillMsgMessage:

	default:
		t.Errorf("msg should be *willMsgMessage")
	}
}

func Test_NewRegisterMessage(t *testing.T) {
	rm := NewMessage(REGISTER)
	switch rm.(type) {
	case *RegisterMessage:

	default:
		t.Errorf("msg should be *registerMessage")
	}
}

func Test_NewRegackMessage(t *testing.T) {
	rm := NewMessage(REGACK)
	switch rm.(type) {
	case *RegackMessage:

	default:
		t.Errorf("msg should be *regackMessage")
	}
}

func Test_NewPublishMessage(t *testing.T) {
	pm := NewMessage(PUBLISH)
	switch pm.(type) {
	case *PublishMessage:

	default:
		t.Errorf("msg should be *publishMessage")
	}
}

func Test_NewPubackMessage(t *testing.T) {
	pm := NewMessage(PUBACK)
	switch pm.(type) {
	case *PubackMessage:

	default:
		t.Errorf("msg should be *pubackMessage")
	}
}

func Test_NewPubcompMessage(t *testing.T) {
	pm := NewMessage(PUBCOMP)
	switch pm.(type) {
	case *PubcompMessage:

	default:
		t.Errorf("msg should be *pubcompMessage")
	}
}

func Test_NewPubrecMessage(t *testing.T) {
	pm := NewMessage(PUBREC)
	switch pm.(type) {
	case *PubrecMessage:

	default:
		t.Errorf("msg should be *pubrecMessage")
	}
}

func Test_NewPubrelMessage(t *testing.T) {
	pm := NewMessage(PUBREL)
	switch pm.(type) {
	case *PubrelMessage:

	default:
		t.Errorf("msg should be *pubrelMessage")
	}
}

func Test_NewSubscribeMessage(t *testing.T) {
	sm := NewMessage(SUBSCRIBE)
	switch sm.(type) {
	case *SubscribeMessage:

	default:
		t.Errorf("msg should be *subscribeMessage")
	}
}

func Test_NewSubackMessage(t *testing.T) {
	sm := NewMessage(SUBACK)
	switch sm.(type) {
	case *SubackMessage:

	default:
		t.Errorf("msg should be *subackMessage")
	}
}

func Test_NewUnsubscribeMessage(t *testing.T) {
	sm := NewMessage(UNSUBSCRIBE)
	switch sm.(type) {
	case *UnsubscribeMessage:

	default:
		t.Errorf("msg should be *unsubscribeMessage")
	}
}

func Test_NewUnsubackMessage(t *testing.T) {
	sm := NewMessage(UNSUBACK)
	switch sm.(type) {
	case *UnsubackMessage:

	default:
		t.Errorf("msg should be *unsubackMessage")
	}
}

func Test_NewPingreqMessage(t *testing.T) {
	pm := NewMessage(PINGREQ)
	switch pm.(type) {
	case *PingreqMessage:

	default:
		t.Errorf("msg should be *pingreqMessage")
	}
}

func Test_NewPingrespMessage(t *testing.T) {
	pm := NewMessage(PINGRESP)
	switch pm.(type) {
	case *PingrespMessage:

	default:
		t.Errorf("msg should be *pingrespMessage")
	}
}

func Test_NewDisconnectMessage(t *testing.T) {
	dm := NewMessage(DISCONNECT)
	switch dm.(type) {
	case *DisconnectMessage:

	default:
		t.Errorf("msg should be *disconnectMessage")
	}
}

func Test_NewWillTopicUpdateMessage(t *testing.T) {
	wm := NewMessage(WILLTOPICUPD)
	switch wm.(type) {
	case *WillTopicUpdateMessage:

	default:
		t.Errorf("msg should be *willTopicUpdateMessage")
	}
}

func Test_NewWillTopicRespMessage(t *testing.T) {
	wm := NewMessage(WILLTOPICRESP)
	switch wm.(type) {
	case *WillTopicRespMessage:

	default:
		t.Errorf("msg should be *willTopicRespMessage")
	}
}

func Test_NewWillMsgUpdateMessage(t *testing.T) {
	wm := NewMessage(WILLMSGUPD)
	switch wm.(type) {
	case *WillMsgUpdateMessage:

	default:
		t.Errorf("msg should be *willMsgUpdateMessage")
	}
}

func Test_NewWillMsgRespMessage(t *testing.T) {
	wm := NewMessage(WILLMSGRESP)
	switch wm.(type) {
	case *WillMsgRespMessage:

	default:
		t.Errorf("msg should be *willMsgRespMessage")
	}
}
