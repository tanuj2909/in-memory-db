package cmd

import "net"

func Ping(conn net.Conn, args []string) {
	if len(args) != 0 {
		conn.Write([]byte("-ERR wrong number of argumnets\r\n"))
		return
	}
	conn.Write([]byte("+PONG\r\n"))
}
