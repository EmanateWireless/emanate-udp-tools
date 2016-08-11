package util

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
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

// MACAddrToBytes converts the given mac-address string into a 6-byte slice
func MACAddrToBytes(mac string) ([6]byte, error) {
	// remove the delimiter ':' characters
	macHex := strings.Replace(mac, ":", "", -1)

	// if the mac is incorrectly formatted
	if len(macHex) != 12 {
		// return the error now
		return [6]byte{}, fmt.Errorf("Incorrect number of octets in mac-address '%s", mac)
	}

	// convert the hex string to bytes
	bytes, err := hex.DecodeString(macHex)

	// convert the bytes into a fixed array
	macBytes := [6]byte{}
	for i, b := range bytes {
		macBytes[i] = b
	}

	// return the mac-address bytes and whether an error occurred
	return macBytes, err
}

// MACBytesToString converts the given mac-address bytes to a formatted string
func MACBytesToString(bytes [6]byte) string {
	// return the formatted string
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
		bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5])
}
