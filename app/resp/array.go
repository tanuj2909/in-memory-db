package resp

import "fmt"

type array struct{}

func (array) Encode(arr []string) []byte {
	bs := bulkString{}
	byteSlice := []byte("*" + fmt.Sprintf("%d", len(arr)) + "\r\n")

	for _, str := range arr {
		byteSlice = append(byteSlice, bs.Encode(str)...)
	}

	return byteSlice
}

func (array) Decode(b []byte) ([]string, []byte, error) {
	if b[0] != '*' {
		return nil, b, fmt.Errorf("invalid format for array: expected the first byte to be '*', got '%q'", b[0])
	}

	n, b, err := parseLen(b[1:])
	if err != nil {
		return nil, b, fmt.Errorf("invalid format for array: %v", err)
	}

	b, err = parseCRLF(b)
	if err != nil {
		return nil, b, fmt.Errorf("invalid format for array: %v", err)
	}

	arr := make([]string, n)
	bulkString := bulkString{}
	str := ""
	for i := 0; i < n; i++ {
		str, b, err = bulkString.Decode(b)
		if err != nil {
			return nil, b, err
		}
		arr[i] = str
	}

	return arr, b, nil
}
