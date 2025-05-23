package cmd

import (
	"encoding/hex"
	"fmt"
	"net"
)

const (
	emptyRDBFile = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"
)

func Psync(conn net.Conn, replID string, offset int) []byte {
	bytes, err := respHandler.String.Encode(fmt.Sprintf("FULLRESYNC %s %d", replID, offset))

	if err != nil {
		fmt.Println("Failed encoding response", err)
		return respHandler.Error.Encode("ERR processing response")
	}

	emptyRDBFile, err := hex.DecodeString(emptyRDBFile)
	if err != nil {
		fmt.Println("Failed to decode hex", err)
		return respHandler.Error.Encode("ERR processing response")
	}

	msg := []byte("$")
	msg = append(msg, []byte(fmt.Sprintf("%d", len(emptyRDBFile)))...)
	msg = append(msg, []byte("\r\n")...)
	msg = append(msg, emptyRDBFile...)

	bytes = append(bytes, msg...)
	return bytes
}
