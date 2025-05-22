package cmd

import (
	"net"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/resp"
	"github.com/tanuj2909/in-memory-db/app/types"
)

var respHandler = resp.RESPHandler{}

func RunCommand(args []string, state *types.ServerState, conn net.Conn) []byte {
	switch strings.ToUpper(args[0]) {
	case "PING":
		return Ping()
	case "ECHO":
		return Echo(args[1])
	case "SET":
		return Set(args[1:]...)
	case "GET":
		return Get(args[1])
	case "INFO":
		return Info(state)
	case "REPLCONF":
		return ReplConf(args[1:], state, conn)
	case "PSYNC":
		return Psync(conn, state.MasterReplId, state.MasterReplOffset)
	}

	return respHandler.Error.Encode("ERR unknown command\r\n")
}
