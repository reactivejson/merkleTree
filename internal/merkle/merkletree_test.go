package merkletree_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	merkletree "github.com/reactivejson/merkleTree/internal/merkle"
	"github.com/reactivejson/merkleTree/internal/merkle/hash"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

var blake3 = hash.NewBlake3()

// stringToByte  turn a string in to a byte array
func stringToByte(input string) []byte {
	x, err := hex.DecodeString(strings.TrimPrefix(input, "0x"))
	if err != nil {
		panic(err)
	}
	return x
}

var tests = []struct {
	// hash type to use
	hashType hash.HashType
	// input to create the node
	data [][]byte
	// expected error when attempting to create the tree
	createErr error
	// root hash after the tree has been created
	root []byte
}{
	{ // 1
		hashType:  blake3,
		data:      [][]byte{},
		createErr: errors.New("the merkle tree should contains at least 1 piece of input"),
	},
	{ // 4
		hashType: blake3,
		data: [][]byte{
			[]byte("Foo"),
		},
		root: stringToByte("941d7de5dacac20e531e260b762b3d3e5a3e13cde2640a88425ca116fc25dcae"),
	},
	{ // 5
		hashType: blake3,
		data: [][]byte{
			[]byte("Foo"),
			[]byte("Bar"),
			[]byte("Baz"),
		},
		root: stringToByte("98c7591c00c07329581ebbeb7acd781c44623f3c3b52405f51cc8028093fd439"),
	},
	{ // 6
		hashType: blake3,
		data: [][]byte{
			[]byte("Foo"),
			[]byte("Bar"),
			[]byte("Baz"),
			[]byte("Qux"),
			[]byte("Quux"),
			[]byte("Quuz"),
		},
		root: stringToByte("906e3a21b148fd4e8150c686eeb505ee7c93eac53cbc9e766546d0567ed75bcd"),
	},
}

func TestNew(t *testing.T) {
	for i, test := range tests {
		tree, err := merkletree.NewTree(test.data, test.hashType)
		if test.createErr != nil {
			assert.Equal(t, test.createErr.Error(), err.Error(), fmt.Sprintf("expected error at test %d", i))
		} else {
			assert.Nil(t, err, fmt.Sprintf("failed to create tree at test %d", i))
			assert.Equal(t, test.root, tree.MerkleRoot(), fmt.Sprintf("unexpected root at test %d", i))
		}

	}
}

func TestProof(t *testing.T) {
	for i, test := range tests {
		if test.createErr == nil {
			tree, err := merkletree.NewTree(test.data, test.hashType)
			assert.Nil(t, err, fmt.Sprintf("failed to create tree at test %d", i))
			for j, data := range test.data {
				proof, err := tree.GenerateMProof(data)
				assert.Nil(t, err, fmt.Sprintf("failed to create proof at test %d input %d", i, j))
				proven, err := merkletree.VerifyMProof(data, proof, tree.MerkleRoot(), blake3)
				assert.Nil(t, err, fmt.Sprintf("error verifying proof at test %d", i))
				assert.True(t, proven, fmt.Sprintf("failed to verify proof at test %d input %d", i, j))
			}
		}
	}
}
func TestMerkleTree_UpdateLeaf(t *testing.T) {
	data := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("merkle"),
	}

	tree, err := merkletree.NewTree(data, blake3)
	assert.NoError(t, err)

	merkleName := []byte("merkle")
	proof, _ := tree.GenerateMProof(merkleName)

	assert.Equal(t, 2, int(proof.Index))
	assert.Equal(t, stringToByte("d3f14150805edec6ae6c7495f92389abe32d4cef58bc4fe279aa5dec75b33f38"), tree.MerkleRoot())
	verified, err := merkletree.VerifyMProof(merkleName, proof, tree.MerkleRoot(), blake3)
	assert.NoError(t, err)
	assert.True(t, verified)

	newLeaf := []byte("newleaf")
	err = tree.UpdateLeaf(2, newLeaf)
	assert.NoError(t, err)

	assert.Equal(t, stringToByte("f192302935da8b624f5c5dbe4c1ae14e50c363bc5933e4a88624d4400d08716f"), tree.MerkleRoot())

	proof, err = tree.GenerateMProof(newLeaf)
	assert.NoError(t, err)

	verified, err = merkletree.VerifyMProof(merkleName, proof, tree.MerkleRoot(), blake3)
	assert.NoError(t, err)
	assert.False(t, verified)

	proof, err = tree.GenerateMProof(newLeaf)
	assert.NoError(t, err)

	//graph = tree.Visual(new(merkletree.StringFormatter), nil)
	//write("visual/tree2.visual", graph)
	//write("visual/proof2.visual", tree.VisualProof(proof, new(merkletree.StringFormatter), nil))

	assert.Equal(t, 2, int(proof.Index))
	ok, err := merkletree.VerifyMProof(newLeaf, proof, tree.MerkleRoot(), blake3)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func write(f, s string) {
	err := os.WriteFile(f, []byte(s), 0644)
	if err != nil {
		panic(err)
	}
}
