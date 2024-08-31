package types

import (
	"golang.org/x/crypto/sha3"
)

func keccak256(data ...[]byte) Hash {
	h := sha3.NewLegacyKeccak256()
	for _, i := range data {
		h.Write(i)
	}
	return MustHashFromBytes(h.Sum(nil), PadNone)
}

func ptr[T any](x T) *T {
	return &x
}
