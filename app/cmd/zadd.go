package cmd

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/tanuj2909/in-memory-db/app/types"
)

func ZADD(server *types.ServerState, args ...string) []byte {
	if len(args) < 3 || (len(args)-1)%2 != 0 {
		return respHandler.Error.Encode("ERR wrong number of arguments for 'zadd' command")
	}

	key := args[0]
	added := 0

	if _, ok := server.SortedSets[key]; !ok {
		server.SortedSets[key] = []types.SortedSetEntry{}
	}

	for i := 1; i < len(args); i += 2 {
		score, err := strconv.ParseFloat(args[i], 64)
		if err != nil {
			return respHandler.Error.Encode(fmt.Sprintf("ERR value is not a valid float: %s", args[i]))
		}
		member := args[i+1]

		found := false
		for j, entry := range server.SortedSets[key] {
			if entry.Member == member {
				server.SortedSets[key][j].Score = score
				found = true
				break
			}
		}
		if !found {
			server.SortedSets[key] = append(server.SortedSets[key], types.SortedSetEntry{
				Member: member,
				Score:  score,
			})
			added++
		}
	}

	sort.Slice(server.SortedSets[key], func(i, j int) bool {
		a, b := server.SortedSets[key][i], server.SortedSets[key][j]
		if a.Score != b.Score {
			return a.Score < b.Score
		}
		return a.Member < b.Member
	})

	return respHandler.Integer.Encode(added)
}
