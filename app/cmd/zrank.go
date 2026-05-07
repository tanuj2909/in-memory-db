package cmd

import "github.com/tanuj2909/in-memory-db/app/types"

func ZRANK(server *types.ServerState, args ...string) []byte {
	if len(args) != 2 {
		return respHandler.Error.Encode("ERR wrong number of arguments for 'zrank' command")
	}

	key, member := args[0], args[1]

	entries, ok := server.SortedSets[key]
	if !ok {
		return respHandler.Null.Encode()
	}

	for i, entry := range entries {
		if entry.Member == member {
			return respHandler.Integer.Encode(i)
		}
	}

	return respHandler.Null.Encode()
}
