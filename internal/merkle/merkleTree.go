package merkletree

import (
	"bytes"
	"errors"
	hash2 "github.com/reactivejson/merkleTree/internal/merkle/hash"
	"math"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

// MerkleTree is the structure for the Merkle tree.
type MerkleTree struct {
	// hash is a pointer to the hashing struct
	hash hash2.HashType
	// data is the data from which the Merkle tree is created
	data [][]byte
	// nodes are the leaf and branch nodes of the Merkle tree
	nodes [][]byte
}

// NewTree creates a new Merkle tree using the provided raw input and default hash type.
// data must contain at least one element for it to be valid.
func NewTree(data [][]byte, hash hash2.HashType) (*MerkleTree, error) {

	if len(data) == 0 {
		return nil, errors.New("the merkle tree should contains at least 1 piece of input")
	}
	if hash == nil {
		return nil, errors.New("please specify hash algo")
	}

	// starts by calculating the number of branches that the tree will have.
	//This is done by finding the next power of 2 greater than or equal to the number of input elements, using the ceil of the log2 of the input length.
	branchesLen := int(math.Exp2(math.Ceil(math.Log2(float64(len(data))))))

	// We pad our input length up to the power of 2.
	nodes := make([][]byte, branchesLen+len(data)+(branchesLen-len(data)))

	// We put the leaves after the branches in the slice of nodes.
	createLeaves(
		data,
		nodes[branchesLen:branchesLen+len(data)],
		hash,
	)
	// Pad the space left after the leaves.
	for i := len(data) + branchesLen; i < len(nodes); i++ {
		nodes[i] = make([]byte, hash.HashLength())
	}

	// Branches.
	createNonLeaves(
		nodes,
		hash,
		branchesLen,
	)

	tree := &MerkleTree{
		hash:  hash,
		nodes: nodes,
		data:  data,
	}

	return tree, nil
}

// Hashes the input slice, placing the result hashes into dest.
func createLeaves(data [][]byte, dest [][]byte, hash hash2.HashType) {
	for i := range data {
		dest[i] = hash.Hash(data[i])
	}
}

// Create the non-leaf nodes from the existing leaf input.
// The function then calls the createNonLeaves function, passing in the nodes slice, the hash function, and the number of branches.
// This function creates the non-leaf nodes of the tree by computing the hash of each pair of child nodes and storing
// it in the corresponding parent node in the slice of nodes.
// The process continues recursively until there is only one node left, which represents the root of the tree.
func createNonLeaves(nodes [][]byte, hash hash2.HashType, leafOffset int) {
	//  iterates through the nodes from the last leaf node to the root node.
	for i := leafOffset - 1; i > 0; i-- {
		// For each non-leaf node, it retrieves the left and right child nodes by accessing the nodes slice with the formula i2 and i2+1, respectively.
		left := nodes[i*2]
		right := nodes[i*2+1]

		// computes the hash of the concatenation of the left and right child nodes
		nodes[i] = hash.Hash(left, right)

	}
}

// GenerateMProof generates the proof for a piece of input.
// If the input is not present in the tree this will return an error.
// If the input is present in the tree this will return the hashes for each level in the tree and the index of the value in the tree.
func (t *MerkleTree) GenerateMProof(data []byte) (*MerkleProof, error) {
	// Find the index of the input
	index, err := t.dataIndex(data)
	if err != nil {
		return nil, err
	}

	// calculates the length of the proof by computing the number of levels required to reach the root of the tree
	proofLen := int(math.Ceil(math.Log2(float64(len(t.data)))))

	//  It initializes an empty slice to hold the hashes of the proof.
	hashes := make([][]byte, proofLen)

	// It initializes the current index to 0.
	currentIndex := 0

	minIndex := uint64(1)

	//  It starts iterating from the index of the input plus half of the length of the nodes slice.
	// The iteration continues until the index reaches the minimum index required for generating the proof.
	//At each iteration, the code computes the sibling hash of the current node and stores it in the hashes slice.
	// The ^1 operation is a bitwise XOR which toggles the last bit of the index, which selects the sibling node in the tree.
	for i := index + uint64(len(t.nodes)/2); i > minIndex; i /= 2 {
		//  stores the computed sibling hash in the hashes slice.
		hashes[currentIndex] = t.nodes[i^1]
		currentIndex++
	}
	return NewProof(hashes, index), nil
}

// MerkleRoot returns the Merkle root (hash of the root node) of the tree.
func (t *MerkleTree) MerkleRoot() []byte {
	// The first element in the slice is not used, and the second element represents the root node of the tree.
	return t.nodes[1]
}

// UpdateLeaf updates the leaf at the specified index with the new input and recalculates the Merkle tree.
func (t *MerkleTree) UpdateLeaf(index uint64, newData []byte) error {

	// Check if index is within bounds.
	if index >= uint64(len(t.data)) {
		return errors.New("index out of bounds")
	}

	// Hash the new input.
	newLeaf := t.hash.Hash(newData)

	// Replace old input with new input.
	t.data[index] = newData

	// Update nodes in the path from the updated leaf to the root.
	nodeIndex := index + uint64(len(t.nodes)/2)
	t.nodes[nodeIndex] = newLeaf
	// Loop through the path from the updated leaf to the root.
	for nodeIndex > 1 {
		// Calculate the index of the sibling node.
		// the ^ operator is the bitwise XOR operator. The line siblingIndex := nodeIndex ^ 1 calculates the index of the sibling node of the current node in the Merkle tree.
		//
		//The ^ operator performs a bitwise XOR operation on the binary representation of the two operands.
		//When the operator is used with the value 1, it flips the least significant bit of the operand,
		//effectively changing the parity of the number (i.e., from even to odd or from odd to even).
		siblingIndex := nodeIndex ^ 1

		// Calculate the index of the parent node.
		parentIndex := nodeIndex / 2

		// Determine if the current node is the left or right child of its parent.
		if nodeIndex%2 == 0 {
			// If it is the left child, calculate the hash of the parent node by hashing
			// the current node's hash and its sibling's hash.
			t.nodes[parentIndex] = t.hash.Hash(t.nodes[nodeIndex], t.nodes[siblingIndex])
		} else {
			// If it is the right child, calculate the hash of the parent node by hashing
			// its sibling's hash and the current node's hash.
			t.nodes[parentIndex] = t.hash.Hash(t.nodes[siblingIndex], t.nodes[nodeIndex])
		}

		nodeIndex = parentIndex
	}

	return nil
}

// dataIndex returns Index of the data in the MerkleTree.
func (t *MerkleTree) dataIndex(input []byte) (uint64, error) {
	for i, data := range t.data {
		if bytes.Equal(data, input) {
			return uint64(i), nil
		}
	}
	return 0, errors.New("data not found")
}
