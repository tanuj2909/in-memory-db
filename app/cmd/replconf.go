package cmd

import (
	"net"

	"github.com/tanuj2909/in-memory-db/app/types"
)

func ReplConf(args []string, serverState *types.ServerState, conn net.Conn) []byte {
	if len(args) == 2 && args[0] == "listening-port" {
		serverState.Replicas = append(
			serverState.Replicas,
			types.Replica{
				Conn: conn,
			},
		)
		res, _ := respHandler.String.Encode("OK")
		return res
	}

	if len(args) == 2 && args[0] == "capa" {
		res, _ := respHandler.String.Encode("OK")
		return res
	}

	return respHandler.Error.Encode("ERR invalid arguments")
}
