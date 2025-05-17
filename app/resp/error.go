package resp

import "fmt"

type errorValue struct{}

func (errorValue) Encode(err string) []byte {
	return []byte("-" + err + "\r\n")
}

func (errorValue) Decode(b []byte) (string, error) {
	if b[0] != '-' {
		return "", fmt.Errorf("invalid format for error expected first character to be '-', got %q", b[0])
	}

	l := len(b)

	if string(b[l-2:]) != "\r\n" {
		return "", fmt.Errorf("empty error value")
	}

	return string(b[1 : l-2]), nil
}
