package _idea

import "crypto/sha256"

const _hashlength = 32

// BLAKE2b is the Blake2b hashing method.
type BLAKE2b struct{}

// New creates a new Blake2b hashing method.
func NewN2b() *BLAKE2b {
	return &BLAKE2b{}
}

// HashLength returns the length of hashes generated by Hash() in bytes.
func (h *BLAKE2b) HashLength() int {
	return _hashlength
}

// Hash generates a BLAKE2b hash from input byte arrays.
func (h *BLAKE2b) Hash(data ...[]byte) []byte {
	var hash [_hashlength]byte
	if len(data) == 1 {
		hash = sha256.Sum256(data[0])
	} else {
		concatDataLen := 0
		for _, d := range data {
			concatDataLen += len(d)
		}
		concatData := make([]byte, concatDataLen)
		curOffset := 0
		for _, d := range data {
			copy(concatData[curOffset:], d)
			curOffset += len(d)
		}
		hash = sha256.Sum256(concatData)
	}

	return hash[:]
}
