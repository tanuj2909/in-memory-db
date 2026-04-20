package cmd

import "github.com/tanuj2909/in-memory-db/app/store"

func Type(args ...string) []byte {
	if len(args) != 1 {
		return respHandler.Error.Encode("ERR wrong number of arguments")
	}
	_, ok := store.Store.Get(args[0])

	var res []byte
	if !ok {
		res, _ = respHandler.String.Encode("none")
	} else {
		res, _ = respHandler.String.Encode("string")
	}

	return res
}
