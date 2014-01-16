package gateway

func chkerr(e error) {
	if e != nil {
		panic(e)
	}
}
