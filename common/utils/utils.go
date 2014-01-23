package utils

import (
	"encoding/binary"
	"fmt"
)

func Bytes2str(bytes []byte) string {
	s := ""
	for _, b := range bytes {
		s += fmt.Sprintf("%X", b)
	}
	return s
}

func B2u16(num []byte) uint16 {
	return binary.BigEndian.Uint16(num)
}

func U162b(num uint16) []byte {
	encNum := make([]byte, 2)
	binary.BigEndian.PutUint16(encNum, num)
	return encNum
}
