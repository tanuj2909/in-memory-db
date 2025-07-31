package types

import (
	"fmt"
	"net"

	"github.com/tanuj2909/in-memory-db/app/resp"
)

type Replica struct {
	Conn              net.Conn
	BytesAcknowledged int
}

func (r *Replica) GetAcknowlegment() error {
	respHandler := resp.RESPHandler{}

	message := respHandler.Array.Encode([]string{"REPLCONF", "GETACK", "*"})
	_, err := r.Conn.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write GETACK command to replica connection: %v", err)
	}
	return nil
}
