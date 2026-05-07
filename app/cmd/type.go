package cmd

import (
	"github.com/tanuj2909/in-memory-db/app/store"
	"github.com/tanuj2909/in-memory-db/app/types"
)

func Type(server *types.ServerState, args ...string) []byte {
	if len(args) != 1 {
		return respHandler.Error.Encode("ERR wrong number of arguments")
	}

	key := args[0]

	if _, ok := store.Store.Get(key); ok {
		res, _ := respHandler.String.Encode("string")
		return res
	}
	if _, ok := server.Streams[key]; ok {
		res, _ := respHandler.String.Encode("stream")
		return res
	}
	if _, ok := server.SortedSets[key]; ok {
		res, _ := respHandler.String.Encode("zset")
		return res
	}

	res, _ := respHandler.String.Encode("none")
	return res
}
