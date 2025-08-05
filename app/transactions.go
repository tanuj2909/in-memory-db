package main

import (
	"net"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/resp"
	"github.com/tanuj2909/in-memory-db/app/types"
	"github.com/tanuj2909/in-memory-db/app/util"
)

func HandleTransactionCommand(conn net.Conn, command string, server *types.ServerState) {

	respHandler := resp.RESPHandler{}
	switch strings.ToUpper(command) {
	case "MULTI":
		if util.IsTransactionStarted(conn, server) {
			util.SendResponse(
				conn,
				respHandler.Error.Encode("ERR MULTI calls can not be nested"),
			)
			util.EndTransaction(conn, server)
			return
		}
		util.StartTransaction(conn, server)
		ok, _ := respHandler.String.Encode("OK")
		util.SendResponse(conn, ok)
	case "EXEC":
		if !util.IsTransactionStarted(conn, server) {
			util.SendResponse(
				conn,
				respHandler.Error.Encode("ERR EXEC without MULTI"),
			)
			return
		}

		transaction, ok := server.Transactions[conn]

		if !ok {
			panic("Cannot find transaction")
		}

		res := util.ExecuteTransaction(&transaction, server)

		if res != nil {
			util.EndTransaction(conn, server)
			util.SendResponse(conn, res)
		}
	case "DISCARD":
		if !util.IsTransactionStarted(conn, server) {
			util.SendResponse(
				conn,
				respHandler.Error.Encode("ERR DISCARD without MULTI"),
			)
			return
		}

		util.EndTransaction(conn, server)
		ok, _ := respHandler.String.Encode("OK")
		util.SendResponse(conn, ok)
	}
}
