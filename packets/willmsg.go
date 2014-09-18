package packets

import (
	"io"
)

type WillMsgMessage struct {
	Header
	WillMsg []byte
}

func (wm *WillMsgMessage) MessageType() byte {
	return WILLMSG
}

func (wm *WillMsgMessage) Write(w io.Writer) error {
	wm.Header.Length = uint16(len(wm.WillMsg) + 2)
	packet := wm.Header.pack()
	packet.WriteByte(WILLMSG)
	packet.Write(wm.WillMsg)
	_, err := packet.WriteTo(w)

	return err
}

func (wm *WillMsgMessage) Unpack(b io.Reader) {
	b.Read(wm.WillMsg)
}
