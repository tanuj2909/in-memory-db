package resp

import "fmt"

type bulkString struct{}

func (bulkString) Encode(s string) []byte {
	lenStr := fmt.Sprintf("%d", len(s))
	return []byte("$" + lenStr + "\r\n" + s + "\r\n")
}

// Reads and decodes a single bulk string.
// Returns the string and the remaing byte slice.
func (bulkString) Decode(b []byte) (string, []byte, error) {
	if b[0] != '$' {
		return "", b, fmt.Errorf("expected first character as '$', got '%q'", b[0])
	}

	n, b, err := parseLen(b[1:])
	if err != nil {
		return "", b, fmt.Errorf("invalid format for bulk string: %v", err)
	}

	//parse starting CR LF
	b, err = parseCRLF(b)
	if err != nil {
		return "", b, fmt.Errorf("invalid format for bulk string: %v", err)
	}

	if len(b) < n+2 {
		return "", b, fmt.Errorf("invalid format for bulk string expected length %d(atleast), got %d", n+2, len(b))
	}

	str := string(b[:n])
	b = b[n:]

	//to parse ending CR LF
	b, err = parseCRLF(b)
	if err != nil {
		return "", b, fmt.Errorf("invalid format for bulk string: %v", err)
	}

	return str, b, nil
}
