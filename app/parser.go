package main

import (
	"fmt"
	"net"

	"github.com/tanuj2909/in-memory-db/app/cmd"
	"github.com/tanuj2909/in-memory-db/app/resp"
	"github.com/tanuj2909/in-memory-db/app/types"
)

func handleCommand(buf []byte, conn net.Conn, state *types.ServerState) {
	respHandler := resp.RESPHandler{}

	arr, next, err := respHandler.DecodeCommand(buf)
	if err != nil {
		fmt.Printf("Error decoding command: %v\n", err)
		return
	}
	res := cmd.RunCommand(arr, state)

	_, err = conn.Write(res)

	if err != nil {
		fmt.Printf("Error writing response to client: %v\n", err)
	}

	if len(next) > 0 {
		handleCommand(next, conn, state)
	}
}
