package main

import (
	"fmt"
	"net"
	"os"
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

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("error reading from connection: ", err)
			}
			return
		}

		handleCommand(buf[:n], conn)
	}
}
