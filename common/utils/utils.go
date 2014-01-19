package utils

import (
	"fmt"
)

func Bytes2str(bytes []byte) string {
	s := ""
	for _, b := range bytes {
		s += fmt.Sprintf("%X", b)
	}
	return s
}
