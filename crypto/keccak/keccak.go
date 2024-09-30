package keccak

import (
	"golang.org/x/crypto/sha3"
)

// Hash is a 32-byte Keccak256 hash.
//
// For most use cases, the Hash type from the types package should be used instead.
type Hash [32]byte

// Hash256 calculates the Keccak256 hash of the given data.
func Hash256(data ...[]byte) (h Hash) {
	k := sha3.NewLegacyKeccak256()
	for _, i := range data {
		k.Write(i)
	}
	copy(h[:], k.Sum(nil))
	return
}
