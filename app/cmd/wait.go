package cmd

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/tanuj2909/in-memory-db/app/types"
)

func Wait(conn net.Conn, state *types.ServerState, args ...string) []byte {
	if len(args) != 2 {
		fmt.Printf("Invalid number of arguments in wait command: expected 2, got %d\n", len(args))
		return nil
	}

	numReplicas, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Error converting number of replicas to integer: %v\n", err)
		return nil
	}

	timeOut, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Error converting timeout to integer: %v\n", err)
		return nil
	}

	for _, replica := range state.Replicas {
		go func(r types.Replica) {
			err := r.GetAcknowlegment()
			if err != nil {
				fmt.Printf("Error getting acknowledgement from replica: %v\n", err)
			}
		}(replica)
	}

	startTime := time.Now()
	for {
		if time.Since(startTime) > time.Duration(timeOut)*time.Millisecond {
			fmt.Printf("Timeout reached waiting for %d replicas\n", numReplicas)
			break
		}

		ackCount := getCorrectAckCount(state.Replicas, state.BytesSent)
		if ackCount >= numReplicas {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	ackCount := getCorrectAckCount(state.Replicas, state.BytesSent)
	return respHandler.Integer.Encode(ackCount)
}

func getCorrectAckCount(replicas []types.Replica, bytesSent int) int {
	ackCount := 0
	for _, replica := range replicas {
		if replica.BytesAcknowledged >= bytesSent {
			ackCount++
		}
	}
	return ackCount
}
