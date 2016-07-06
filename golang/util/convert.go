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
