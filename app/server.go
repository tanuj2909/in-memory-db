package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/types"
	"github.com/tanuj2909/in-memory-db/app/util"
)

func NewServerState(args *Args) *types.ServerState {
	state := types.ServerState{
		Role:             "master",
		Port:             args.port,
		MasterReplId:     util.RandomAlphanumeric(40),
		MasterReplOffset: 0,
	}

	if args.replicaof != "" {
		state.Role = "slave"
		host := strings.Split(args.replicaof, " ")
		state.MasterHost = host[0]
		state.MasterPort = host[1]
		state.MasterReplId = "?"
		state.MasterReplOffset = -1
		handshakeWithMaster(&state)
	}

	return &state
}
func main() {
	args := GetArgs()

	serverState := NewServerState(&args)

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", args.port))

	if err != nil {
		fmt.Printf("Failed to bind port %d\n", args.port)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn, serverState)

	}
}

func handleConnection(conn net.Conn, state *types.ServerState) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("error reading from connection: ", err)
			}
			return
		}

		handleCommand(buf[:n], conn, state)
	}
}
