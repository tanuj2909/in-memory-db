package cmd

import (
	"fmt"
	"net"

	"github.com/tanuj2909/in-memory-db/app/store"
)

func Get(conn net.Conn, args []string) {
	if len(args) != 1 {
		conn.Write([]byte("-ERR wrong number of argumnets\r\n"))
		return
	}
	val, ok := store.Store.Get(args[0])
	if !ok {
		conn.Write([]byte("$-1\r\n"))
		return
	}
	resp := fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
	conn.Write([]byte(resp))
}
