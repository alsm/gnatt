package packets

import (
	"io"
)

type SearchGwMessage struct {
	Header
	Radius byte
}

func (s *SearchGwMessage) MessageType() byte {
	return SEARCHGW
}

func (s *SearchGwMessage) Write(w io.Writer) error {
	packet := s.Header.pack()
	packet.WriteByte(SEARCHGW)
	packet.WriteByte(s.Radius)
	_, err := packet.WriteTo(w)

	return err
}

func (s *SearchGwMessage) Unpack(b io.Reader) {
	s.Radius = readByte(b)
}
