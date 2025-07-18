package cmd

import (
	"strconv"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/store"
)

func Set(isMaster bool, args ...string) []byte {
	shouldReply := !isMaster
	if len(args) != 2 && len(args) != 4 {
		return respHandler.Error.Encode("ERR wrong number of argumnets")
	}
	ttl := int64(0)
	if len(args) == 4 {
		if strings.ToUpper(args[2]) != "EX" {
			return respHandler.Error.Encode("ERR wrong usage of set")
		}
		var err error
		ttl, err = strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return respHandler.Error.Encode("ERR EX argument must be an integer")
		}
	}

	store.Store.Set(args[0], args[1], ttl)
	if shouldReply {
		res, _ := respHandler.String.Encode("OK")
		return res
	}

	return nil
}
