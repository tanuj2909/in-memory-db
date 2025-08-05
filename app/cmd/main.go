package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/resp"
	"github.com/tanuj2909/in-memory-db/app/types"
)

var respHandler = resp.RESPHandler{}

func RunCommand(args []string, state *types.ServerState, conn net.Conn, buf []byte, isMaster bool) []byte {
	var res []byte
	switch strings.ToUpper(args[0]) {
	case "PING":
		res = Ping(isMaster)
	case "ECHO":
		res = Echo(args[1])
	case "SET":
		res = Set(isMaster, args[1:]...)
		if state.Role == "master" {
			state.BytesSent += len(buf)
			streamToReplicas(state.Replicas, buf)
		}
	case "GET":
		res = Get(args[1])
	case "INFO":
		res = Info(state)
	case "REPLCONF":
		res = ReplConf(args[1:], state, conn)
	case "PSYNC":
		res = Psync(conn, state.MasterReplId, state.MasterReplOffset)
	case "WAIT":
		res = Wait(conn, state, args[1:]...)
	case "CONFIG":
		res = Config(state, args[1:]...)
	case "INCR":
		res = Incr(args[1])
	default:
		res = respHandler.Error.Encode("ERR unknown command\r\n")
	}

	return res
}

func streamToReplicas(replicas []types.Replica, buf []byte) {
	for _, replica := range replicas {
		_, err := replica.Conn.Write(buf)
		if err != nil {
			fmt.Printf("falied to write to replica: %v\n", err.Error())
		}
	}
}
