package cmd

import (
	"strings"

	"github.com/tanuj2909/in-memory-db/app/resp"
)

var respHandler = resp.RESPHandler{}

func RunCommand(args []string) []byte {
	switch strings.ToUpper(args[0]) {
	case "PING":
		return Ping()
	case "ECHO":
		return Echo(args[1])
	case "SET":
		return Set(args[1:]...)
	case "GET":
		return Get(args[1])
	}

	return respHandler.Error.Encode("ERR unknown command\r\n")
}
