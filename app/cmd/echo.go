package cmd

import (
	"fmt"
	"net"
)

func Echo(conn net.Conn, args []string) {
	if len(args) != 1 {
		conn.Write([]byte("-ERR wrong number of arguments\r\n"))
		return
	}

	arg := args[0]
	resp := fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
	conn.Write([]byte(resp))
}
