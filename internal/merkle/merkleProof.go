package merkletree

import (
	"bytes"
	"github.com/reactivejson/merkleTree/internal/merkle/hash"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

// MerkleProof is a proof of a Merkle tree.
type MerkleProof struct {
	Hashes [][]byte //2D byte array representing the hashes of nodes in the Merkle tree
	Index  uint64   // The index of the input element for which the proof was generated
}

// NewProof generates a Merkle proof.
func NewProof(hashes [][]byte, index uint64) *MerkleProof {
	return &MerkleProof{
		Hashes: hashes,
		Index:  index,
	}
}

// VerifyMProof verifies a Merkle tree proof for a piece of input using the default hash type.
// The proof and path are as per Merkle tree's GenerateMProof(), and root is the root hash of the tree against which the proof is to
// be verified.  Note that this does not require the Merkle tree to verify the proof, only its root; this allows for checking
// against historical trees without having to instantiate them.
//
// This returns true if the proof is verified, otherwise false.
func VerifyMProof(data []byte, proof *MerkleProof, root []byte, hashType hash.HashType) (bool, error) {
	proofHash := proofHash(data, proof, hashType)
	if bytes.Equal(root, proofHash) {
		// If the hash in the root matches the proof hash, this line returns true and a nil error.
		return true, nil
	}

	//If the proof is not verified, this line returns false and a nil error.
	return false, nil
}

// proofHash generates a proof hash for a piece of input using the provided Merkle proof and hash function.
func proofHash(data []byte, proof *MerkleProof, hashType hash.HashType) []byte {

	var proofHash []byte

	// Generate the initial hash by hashing the input with the provided hash function.
	proofHash = hashType.Hash(data)

	// Calculate the starting index in the proof array based on the number of hashes in the proof.
	index := proof.Index + (1 << uint(len(proof.Hashes)))

	// Loop over each hash in the proof array, combining them with the proof hash based on whether the current index is even or odd.
	for _, hash := range proof.Hashes {
		if index%2 == 0 {
			// If the index is even, hash the proof hash and the current hash together.
			proofHash = hashType.Hash(proofHash, hash)
		} else {
			// If the index is odd, hash the current hash and the proof hash together.
			proofHash = hashType.Hash(hash, proofHash)
		}
		// Shift the index right by one bit, effectively dividing it by 2 and rounding down to the nearest integer.
		index >>= 1
	}

	// Return the final proof hash.
	return proofHash
}
