package util

import (
	"crypto/rand"
	"math/big"
)

const charSet = "qwertyuiopasdfghjklzxcvbnnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func RandomAlphanumeric(n int) string {
	res := make([]byte, n)

	for i := range res {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		}
		res[i] = charSet[num.Int64()]
	}
	return string(res)
}
