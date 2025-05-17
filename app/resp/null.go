package resp

type null struct{}

const NULL_VALUE = "$-1\r\n"

func (null) Encode() []byte {
	return []byte(NULL_VALUE)
}
