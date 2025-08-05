package main

import (
	"fmt"
	"net"

	"github.com/tanuj2909/in-memory-db/app/cmd"
	"github.com/tanuj2909/in-memory-db/app/resp"
	"github.com/tanuj2909/in-memory-db/app/types"
	"github.com/tanuj2909/in-memory-db/app/util"
)

func handleCommand(buf []byte, conn net.Conn, state *types.ServerState, isMaster bool) {
	respHandler := resp.RESPHandler{}

	arr, next, err := respHandler.DecodeCommand(buf)
	if err != nil {
		fmt.Printf("Error decoding command: %v\n", err)
		return
	}
	if util.IsTransactionCommand(arr[0]) {
		HandleTransactionCommand(conn, arr[0], state)
	} else {
		inTransaction := util.IsTransactionStarted(conn, state)

		if inTransaction {
			state.TransactionMutex.Lock()
			state.Transactions[conn] = types.TransactionData{
				Started: true,
				Queue:   append(state.Transactions[conn].Queue, buf),
			}
			state.TransactionMutex.Unlock()

			res, err := respHandler.String.Encode("QUEUED")
			if err != nil {
				fmt.Printf("Error encoding response: %v\n", err)
			}
			_, err = conn.Write(res)
			if err != nil {
				fmt.Printf("Error writing response to client: %v\n", err)
			}
		} else {
			res := cmd.RunCommand(arr, state, conn, buf, isMaster)

			if res != nil {
				_, err := conn.Write(res)
				if err != nil {
					fmt.Printf("Error writing response to client: %v\n", err)
				}
			}
		}
	}

	if isMaster {
		state.AckOffset += len(buf)
	}

	if len(next) > 0 {
		handleCommand(next, conn, state, isMaster)
	}
}
