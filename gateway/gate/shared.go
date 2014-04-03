package gateway

func validateClientId(clientid []byte) (string, error) {
	if len(clientid) == 0 {
		ERROR.Println("zero length client id not allowed")
		return "", ErrZeroLengthClientID
	}
	if len(clientid) > 23 {
		ERROR.Println("client id longer than 23 characters")
		return "", ErrClientIDTooLong
	}
	return string(clientid), nil
}
