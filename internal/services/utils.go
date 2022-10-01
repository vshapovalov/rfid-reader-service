package services

import "math/big"

func byteArrayToInt(b []byte) int {
	return int(big.NewInt(0).SetBytes(b).Uint64())
}

func reverseArray[T comparable](arr []T) []T {
	result := make([]T, len(arr))
	for i, v := range arr {
		result[len(arr)-1-i] = v
	}
	return result
}
