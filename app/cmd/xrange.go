package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tanuj2909/in-memory-db/app/types"
)

func Xrange(server *types.ServerState, args ...string) []byte {
	if len(args) < 3 {
		return respHandler.Error.Encode(
			fmt.Sprintf("ERR wrong number of arguments for '%s' command, expected at least 3, for %d", args[0], len(args)),
		)
	}

	streamKey := args[0]
	stream, ok := server.Streams[streamKey]

	if !ok {
		return respHandler.Array.Encode([]string{})
	}

	items := fetchFromStream(stream, args[1], args[2])
	encodedItems, err := encodeStreamEntrySlice(items)
	if err != nil {
		fmt.Printf("Error encoding stream entries: %s\n", err)
		return nil
	}
	return encodedItems
}

func fetchFromStream(streams []types.StreamEntry, start string, end string) []types.StreamEntry {
	res := make([]types.StreamEntry, 0)
	for _, stream := range streams {
		streamIDParts := strings.Split(stream.ID, "-")

		streamTimestamp, err := strconv.Atoi(streamIDParts[0])
		if err != nil {
			panic(fmt.Sprintf("Invalid ID format present: %s", stream.ID))
		}
		streamSeqNo, err := strconv.Atoi(streamIDParts[1])
		if err != nil {
			panic(fmt.Sprintf("Invalid ID format present: %s", stream.ID))
		}

		if isGreaterOrEqual(streamTimestamp, streamSeqNo, start) && isLessThanOrEqual(streamTimestamp, streamSeqNo, end) {
			res = append(res, stream)
		}
	}

	return res
}

func isGreaterOrEqual(streamTimestamp int, streamSeqNo int, startID string) bool {
	if startID == "-" {
		return true
	}

	idParts := strings.Split(startID, "-")
	idTimestamp, err := strconv.Atoi(idParts[0])
	if err != nil {
		panic(fmt.Sprintf("Invalid ID format present: %s", startID))
	}

	if len(idParts) == 1 {
		return streamTimestamp >= idTimestamp
	}

	idSeqNo, err := strconv.Atoi(idParts[1])
	if err != nil {
		panic(fmt.Sprintf("Invalid ID format present: %s", startID))
	}

	return streamTimestamp > idTimestamp || (streamTimestamp == idTimestamp && streamSeqNo >= idSeqNo)
}

func isLessThanOrEqual(streamTimestamp int, streamSeqNo int, endID string) bool {
	if endID == "+" {
		return true
	}

	idParts := strings.Split(endID, "-")
	idTimestamp, err := strconv.Atoi(idParts[0])
	if err != nil {
		panic(fmt.Sprintf("Invalid ID format present: %s", endID))
	}

	if len(idParts) == 1 {
		return streamTimestamp <= idTimestamp
	}

	idSeqNo, err := strconv.Atoi(idParts[1])
	if err != nil {
		panic(fmt.Sprintf("Invalid ID format present: %s", endID))
	}

	return streamTimestamp < idTimestamp || (streamTimestamp == idTimestamp && streamSeqNo <= idSeqNo)
}

func encodeStreamEntrySlice(entries []types.StreamEntry) ([]byte, error) {
	encodedBytes := []byte(fmt.Sprintf("*%d\r\n", len(entries)))

	for _, entry := range entries {
		encodedEntry, err := encodeStreamEntry(entry)
		if err != nil {
			return nil, fmt.Errorf("failed to encode entry: %v", err)
		}

		encodedBytes = append(encodedBytes, encodedEntry...)
	}

	return encodedBytes, nil
}

func encodeStreamEntry(entry types.StreamEntry) ([]byte, error) {
	encodedBytes := []byte(fmt.Sprintf("*%d\r\n", 2))

	encodedID, err := respHandler.String.Encode(entry.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to encode ID: %v", err)
	}
	encodedBytes = append(encodedBytes, encodedID...)

	kvSlice := make([]string, 0)
	for key, value := range entry.KV {
		kvSlice = append(kvSlice, key, value)
	}
	encodedKVs := respHandler.Array.Encode(kvSlice)
	encodedBytes = append(encodedBytes, encodedKVs...)

	return encodedBytes, nil
}
