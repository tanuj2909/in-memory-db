package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/resp"
	"github.com/tanuj2909/in-memory-db/app/types"
)

func handshakeWithMaster(state *types.ServerState) {
	masterConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", state.MasterHost, state.MasterPort))
	if err != nil {
		fmt.Println("Failed to connect to master: ", err)
		return
	}

	respHandler := resp.RESPHandler{}
	//PING
	err = sendAndAssertReply(
		masterConn,
		[]string{"PING"},
		"PONG",
		respHandler,
	)

	if err != nil {
		fmt.Println("Failed to send PING to master: ", err)
		return
	}
	//REPLCONF listening-port <port>
	err = sendAndAssertReply(
		masterConn,
		[]string{"REPLCONF", "listening-port", fmt.Sprintf("%d", state.Port)},
		"OK",
		respHandler,
	)
	if err != nil {
		fmt.Println("Failed to send REPLCONF listening-port to master: ", err)
		return
	}

	//REPLCONF capa psync2
	err = sendAndAssertReply(
		masterConn,
		[]string{"REPLCONF", "capa", "psync2"},
		"OK",
		respHandler,
	)
	if err != nil {
		fmt.Println("Failed to send REPLCONF capa psync2 to master: ", err)
		return
	}

	//PSYNC <replicationid> <offset>
	err = sendAndGetData(
		masterConn,
		[]string{"PSYNC", "?", fmt.Sprintf("%d", -1)},
		respHandler,
		state,
	)

	if err != nil {
		fmt.Println("Failed to send PSYNC to master: ", err)
		return
	}
}

func sendAndGetData(conn net.Conn, msgArr []string, respHandler resp.RESPHandler, state *types.ServerState) error {
	bytes := respHandler.Array.Encode(msgArr)
	conn.Write(bytes)

	resp := make([]byte, 1024)
	n, _ := conn.Read(resp)
	res, remain, err := respHandler.String.Decode(resp[:n])
	if err != nil {
		return fmt.Errorf("failed to decode response: %s", err)
	}
	if len(remain) > 0 {
		return fmt.Errorf("unexpected remaining bytes: %q", remain)
	}

	responseParts := strings.Split(res, " ")
	if len(responseParts) != 3 {
		return fmt.Errorf("expected 3 parts in PSYNC response, got %d", len(responseParts))
	}
	if responseParts[0] != "FULLRESYNC" {
		return fmt.Errorf("expected FULLRESYNC in PSYNC response, got %s", responseParts[0])
	}
	state.MasterReplId = responseParts[1]
	portAsInt, err := strconv.Atoi(responseParts[2])
	if err != nil {
		return fmt.Errorf("failed to convert port to int: %s", err)
	}
	state.MasterReplOffset = portAsInt

	return nil

}

func sendAndAssertReply(conn net.Conn, msgArr []string, expectedMsg string, respHandler resp.RESPHandler) error {
	bytes := respHandler.Array.Encode(msgArr)
	conn.Write(bytes)

	resp := make([]byte, 1024)
	n, _ := conn.Read(resp)
	msg, remain, err := respHandler.String.Decode(resp[:n])
	if err != nil {
		return fmt.Errorf("failed to decode response: %s", err)
	}
	if msg != expectedMsg {
		return fmt.Errorf("expected PONG got %s", string(resp[:n]))
	}
	if len(remain) > 0 {
		return fmt.Errorf("unexpected remaining bytes: %q", remain)
	}

	return nil
}
