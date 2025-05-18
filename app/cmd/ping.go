package cmd

import "fmt"

func Ping() []byte {
	res, err := respHandler.String.Encode("PONG")
	if err != nil {
		fmt.Printf("Error encoding response: %s\n", err)
		return respHandler.Null.Encode()
	}
	return res
}
