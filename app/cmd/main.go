package cmd

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func HandleRequest(conn net.Conn, reader *bufio.Reader, line string) {
	cmnd, args := HandleCommand(reader, line)

	switch strings.ToUpper(cmnd) {
	case "PING":
		Ping(conn, args)
	case "ECHO":
		Echo(conn, args)
	case "SET":
		Set(conn, args)
	case "GET":
		Get(conn, args)
	default:
		conn.Write([]byte("-ERR unknown command\r\n"))
	}
}

func HandleCommand(reader *bufio.Reader, line string) (string, []string) {
	numArgs := 0
	fmt.Sscanf(line, "*%d", &numArgs)

	args := make([]string, 0, numArgs)

	for i := 0; i < numArgs; i++ {
		reader.ReadString('\n')
		arg, _ := reader.ReadString('\n')
		args = append(args, strings.TrimSpace(arg))
	}
	if len(args) == 0 {
		return "", []string{}
	}
	return args[0], args[1:]
}
