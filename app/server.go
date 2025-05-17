package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/cmd"
)

func main() {
	args := GetArgs()

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
		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("error reading from connection: ", err)
			}
			return
		}

		if strings.HasPrefix(line, "*") {
			cmd.HandleRequest(conn, reader, line)
		}
	}
}
