package cmd

func Echo(msg string) []byte {
	return respHandler.BulkString.Encode(msg)
}
