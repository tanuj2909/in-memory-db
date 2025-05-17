package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/store"
)

func Set(conn net.Conn, args []string) {
	if len(args) < 2 {
		conn.Write([]byte("-ERR wrong number of argumnets\r\n"))
		return
	}
	ttl := 0
	if len(args) == 4 && strings.ToUpper(args[2]) == "EX" {
		fmt.Sscanf(args[3], "%d", &ttl)
	}
	store.Store.Set(args[0], args[1], ttl)
	conn.Write([]byte("+OK\r\n"))
}
