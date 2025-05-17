package resp

import (
	"fmt"
	"strconv"
)

type integer struct{}

func (integer) Encode(n int) []byte {
	intResp := fmt.Sprintf(":%d\r\n", n)
	return []byte(intResp)
}

func (integer) Decode(b []byte) (int, error) {
	if b[0] != ':' {
		return 0, fmt.Errorf("invalid format for integer expected first character to be ':', got %q", b[0])
	}

	l := len(b)
	if string(b[l-2:]) != "\r\n" {
		return 0, fmt.Errorf("invalid format must end with CR LF, but ends with %s", string(b[l-2:]))
	}

	n, err := strconv.Atoi(string(b[1 : l-2]))
	if err != nil {
		return 0, fmt.Errorf("can't convert the passed bytes to integer: %v", err)
	}

	return n, err
}
