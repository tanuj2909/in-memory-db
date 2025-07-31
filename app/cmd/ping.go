package cmd

import "fmt"

func Ping(isMasterCommand bool) []byte {
	if isMasterCommand {
		return nil
	}

	res, err := respHandler.String.Encode("PONG")
	if err != nil {
		fmt.Printf("Error encoding response: %s\n", err)
		return respHandler.Null.Encode()
	}
	return res
}
