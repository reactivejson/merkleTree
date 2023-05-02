package hash

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

// HashFunc is a hashing function.
type HashFunc func(...[]byte) []byte

// HashType defines the interface that must be supplied by hash functions.
type HashType interface {
	// Hash calculates the hash of a given input.
	Hash(...[]byte) []byte

	// HashLength provides the length of the hash.
	HashLength() int
}
