package cmd

import (
	"fmt"
	"net"
)

func Psync(conn net.Conn, replID string, offset int) []byte {
	bytes, err := respHandler.String.Encode(fmt.Sprintf("FULLRESYNC %s %d", replID, offset))

	if err != nil {
		return respHandler.Error.Encode("ERR encoding response")
	}

	return bytes
}
