package util

import (
	"fmt"
	"net"

	"github.com/tanuj2909/in-memory-db/app/cmd"
	"github.com/tanuj2909/in-memory-db/app/resp"
	"github.com/tanuj2909/in-memory-db/app/types"
)

func SendResponse(conn net.Conn, response []byte) {
	_, err := conn.Write(response)
	if err != nil {
		fmt.Printf("Error writing response to client: %v\n", err)
	}
}

func StartTransaction(conn net.Conn, server *types.ServerState) {
	server.TransactionMutex.Lock()
	server.Transactions[conn] = types.TransactionData{
		Started: true,
		Queue:   [][]byte{},
	}
	server.TransactionMutex.Unlock()
}

func IsTransactionCommand(command string) bool {
	return command == "MULTI" || command == "EXEC" || command == "DISCARD"
}

func IsTransactionStarted(conn net.Conn, server *types.ServerState) bool {
	server.TransactionMutex.Lock()
	defer server.TransactionMutex.Unlock()

	transaction, ok := server.Transactions[conn]
	return ok && transaction.Started
}

func EndTransaction(conn net.Conn, server *types.ServerState) {
	server.TransactionMutex.Lock()
	defer server.TransactionMutex.Unlock()
	server.Transactions[conn] = types.TransactionData{
		Started: false,
		Queue:   [][]byte{},
	}

}

func ExecuteTransaction(t *types.TransactionData, server *types.ServerState) []byte {
	output := make([][]byte, 0)
	respHandler := resp.RESPHandler{}
	for _, command := range t.Queue {
		arr, _, err := respHandler.DecodeCommand(command)

		if err != nil {
			fmt.Printf("Error decoding queued command")
			return nil
		}

		response := cmd.RunCommand(arr, server, nil, command, false)

		output = append(output, response)
	}

	return respHandler.Array.EncodeFromBytes(output)
}
