package resp

type RESPHandler struct {
	String     simpleString
	BulkString bulkString
	Array      array
	Null       null
	Integer    integer
	Error      errorValue
}

func (h RESPHandler) DecodeCommand(b []byte) ([]string, []byte, error) {
	return h.Array.Decode(b)
}
