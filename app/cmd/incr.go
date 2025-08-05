package cmd

import (
	"github.com/tanuj2909/in-memory-db/app/store"
)

func Incr(key string) []byte {
	val, ok := store.Store.Incr(key)
	if !ok {
		return respHandler.Error.Encode("ERR value is not an integer or out of range")
	}
	return respHandler.Integer.Encode(val)
}
