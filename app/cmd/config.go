package cmd

import (
	"fmt"

	"github.com/tanuj2909/in-memory-db/app/types"
)

func Config(state *types.ServerState, args ...string) []byte {
	if len(args) != 2 {
		return respHandler.Error.Encode(fmt.Sprintf("ERR wrong number of argumnets: expected2, fot %d\n", len(args)))
	}

	if args[0] == "GET" {
		if args[1] == "dir" {
			return respHandler.Array.Encode([]string{"dir", state.DBDir})
		}

		if args[1] == "dbfilename" {
			return respHandler.Array.Encode([]string{"dbfilename", state.DBFileName})
		}
	}

	return respHandler.Error.Encode("ERR unknown command")
}
