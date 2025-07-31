package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/resp"
	"github.com/tanuj2909/in-memory-db/app/types"
)

var respHandler = resp.RESPHandler{}

func RunCommand(args []string, state *types.ServerState, conn net.Conn, buf []byte, isMaster bool) {
	var err error
	switch strings.ToUpper(args[0]) {
	case "PING":
		_, err = conn.Write(Ping(isMaster))
	case "ECHO":
		_, err = conn.Write(Echo(args[1]))
	case "SET":
		_, err = conn.Write(Set(isMaster, args[1:]...))
		if state.Role == "master" {
			state.BytesSent += len(buf)
			streamToReplicas(state.Replicas, buf)
		}
	case "GET":
		_, err = conn.Write(Get(args[1]))
	case "INFO":
		_, err = conn.Write(Info(state))
	case "REPLCONF":
		_, err = conn.Write(ReplConf(args[1:], state, conn))
	case "PSYNC":
		_, err = conn.Write(Psync(conn, state.MasterReplId, state.MasterReplOffset))
	case "WAIT":
		_, err = conn.Write(Wait(conn, state, args[1:]...))
	default:
		_, err = conn.Write(respHandler.Error.Encode("ERR unknown command\r\n"))
	}

	if err != nil {
		fmt.Printf("Error writing response to client: %v\n", err)
	}
}

func streamToReplicas(replicas []types.Replica, buf []byte) {
	for _, replica := range replicas {
		_, err := replica.Conn.Write(buf)
		if err != nil {
			fmt.Printf("falied to write to replica: %v\n", err.Error())
		}
	}
}
