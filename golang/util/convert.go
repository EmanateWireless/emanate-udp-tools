package util

import (
	"encoding/binary"
	"math"
)

// Float32FromBytes creates a float32 from a byte array
func Float32FromBytes(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

// UTF16ToASCII converts the given UTF-16 string to an ASCII string
func UTF16ToASCII(u string) string {
	// create the target byte slice
	a := []byte{}

	// get the utf16 string bytes
	ub := []byte(u)
	ubl := len(ub)

	// iterate through all of the bytes
	for i := 0; (ubl - i) >= 2; i += 2 {
		// skip the first zero-byte and use only the second byte
		a = append(a, ub[i+1])
	}

	// return the ascii string
	return string(a)
}
