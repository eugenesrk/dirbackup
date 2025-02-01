package tools

import (
	"crypto/subtle"
)

func ConstantTimeCompare(inA, inB string) bool {
	if subtle.ConstantTimeEq(int32(len(inA)), int32(len(inB))) != 1 {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(inA), []byte(inB)) != 1 {
		return false
	}

	return true
}
