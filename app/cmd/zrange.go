package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/types"
)

func ZRANGE(server *types.ServerState, args ...string) []byte {
	if len(args) < 3 {
		return respHandler.Error.Encode("ERR wrong number of arguments for 'zrange' command")
	}

	key := args[0]
	withScores := len(args) > 3 && strings.ToUpper(args[3]) == "WITHSCORES"

	entries, ok := server.SortedSets[key]
	if !ok {
		return respHandler.Array.Encode([]string{})
	}

	n := len(entries)

	start, err := strconv.Atoi(args[1])
	if err != nil {
		return respHandler.Error.Encode("ERR value is not an integer or out of range")
	}
	stop, err := strconv.Atoi(args[2])
	if err != nil {
		return respHandler.Error.Encode("ERR value is not an integer or out of range")
	}

	if start < 0 {
		start = n + start
	}
	if stop < 0 {
		stop = n + stop
	}
	if start < 0 {
		start = 0
	}
	if stop >= n {
		stop = n - 1
	}

	if start > stop {
		return respHandler.Array.Encode([]string{})
	}

	slice := entries[start : stop+1]

	if withScores {
		result := make([]string, 0, len(slice)*2)
		for _, entry := range slice {
			result = append(result, entry.Member, fmt.Sprintf("%g", entry.Score))
		}
		return respHandler.Array.Encode(result)
	}

	result := make([]string, 0, len(slice))
	for _, entry := range slice {
		result = append(result, entry.Member)
	}
	return respHandler.Array.Encode(result)
}
