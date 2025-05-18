package cmd

import (
	"github.com/tanuj2909/in-memory-db/app/store"
)

func Get(key string) []byte {
	val, ok := store.Store.Get(key)
	if !ok {
		return respHandler.Null.Encode()
	}

	return respHandler.BulkString.Encode(val)
}
