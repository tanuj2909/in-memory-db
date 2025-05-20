package cmd

import (
	"fmt"

	"github.com/tanuj2909/in-memory-db/app/types"
)

func Info(state *types.ServerState) []byte {
	res := fmt.Sprintf("role:%s", state.Role)
	res += fmt.Sprintf("\nmaster_replid:%s", state.MasterReplId)
	res += fmt.Sprintf("\nmaster_repl_offset:%d", state.MasterReplOffset)
	return respHandler.BulkString.Encode(res)
}
