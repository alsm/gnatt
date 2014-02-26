package gateway

import (
	"fmt"
)

func validateClientId(clientid []byte) (string, error) {
	if len(clientid) == 0 {
		return "", fmt.Errorf("zero length client id not allowed")
	}
	if len(clientid) > 23 {
		return "", fmt.Errorf("client id too long")
	}
	return string(clientid), nil
}
