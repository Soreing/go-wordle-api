package hexcode

import (
	"math/rand"
)

var HEXDEC_CHARSET[16]byte = [16]byte{
	'0','1','2','3','4','5','6','7',
	'8','9','A','B','C','D','E','F',
}

// Generates a Hexadecimal code string of n bytes
func Generate(bytes int) string{
	data := make([]byte, bytes)
	for i := 0; i<bytes; i++{
		data[i] = HEXDEC_CHARSET[rand.Intn(16)]
	}

	return string(data[:])
}