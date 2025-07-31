package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/tanuj2909/in-memory-db/app/types"
)

func ReplConf(args []string, serverState *types.ServerState, conn net.Conn) []byte {

	// REPLCONF listening-port
	if len(args) >= 2 && args[0] == "listening-port" {
		serverState.Replicas = append(
			serverState.Replicas,
			types.Replica{
				Conn:              conn,
				BytesAcknowledged: 0,
			},
		)
		res, _ := respHandler.String.Encode("OK")
		return res
	}

	// REPLCONF capa
	if len(args) >= 2 && args[0] == "capa" {
		res, _ := respHandler.String.Encode("OK")
		return res
	}

	// REPLCONF ACK *
	if len(args) >= 2 && args[0] == "GETACK" && args[1] == "*" {
		return respHandler.Array.Encode([]string{"REPLCONF", "ACK", fmt.Sprintf("%d", serverState.AckOffset)})
	}

	// REPLCONF ACK <bytes>
	if len(args) >= 2 && args[0] == "GETACK" {
		bytesOffset, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Error converting bytes offset to integer: %v\n", err)
			return nil
		}

		for ind, replica := range serverState.Replicas {
			if replica.Conn == conn {
				serverState.Replicas[ind].BytesAcknowledged = bytesOffset
				fmt.Printf("Bytes acknowlefged by replica (%s) updated: %d\n", replica.Conn.RemoteAddr().String(), bytesOffset)
				return nil
			}
		}
		return nil
	}

	return respHandler.Error.Encode("ERR invalid arguments")
}
