package resp

import "fmt"

type simpleString struct{}

func (simpleString) Encode(s string) ([]byte, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' || s[i] == '\r' {
			return nil, fmt.Errorf("invalid character in simple string: %q", s[i])
		}
	}

	return []byte("+" + s + "\r\n"), nil
}

func (simpleString) Decode(b []byte) (string, []byte, error) {
	if b[0] != '+' {
		return "", nil, fmt.Errorf("invalid format for simple string expected first character to be '+', got %q", b[0])
	}

	for i := 1; i < len(b); i++ {
		if b[i] == '\r' {
			if i+1 < len(b) && b[i+1] == '\n' {
				return string(b[1:i]), b[i+2:], nil
			}
			return "", nil, fmt.Errorf("invalid format for simple string: expected the last two bytes to be CR LF, got %q", b[i:])
		}
	}

	return "", nil, fmt.Errorf("invalid format for simple string expected the last two bytes to be CR LF, got %q", b[len(b)-1:])
}
