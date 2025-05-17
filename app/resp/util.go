package resp

import (
	"fmt"
	"strconv"
)

func parseLen(b []byte) (int, []byte, error) {
	lenStr := ""
	for i := 0; i < len(b); i++ {
		if b[i] == '\r' {
			break
		}
		lenStr += string(b[i])
	}

	n, err := strconv.Atoi(lenStr)
	if err != nil {
		return 0, nil, fmt.Errorf("cannot parse number(length) from string %s: %v", lenStr, err)
	}

	return n, b[len(lenStr):], nil
}

func parseCRLF(b []byte) ([]byte, error) {
	if b[0] != '\r' || b[1] != '\n' {
		return nil, fmt.Errorf("expexted the next two bytes to be CR LF, go %q", b[0:2])
	}

	return b[2:], nil
}
