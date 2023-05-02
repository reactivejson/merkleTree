package merkletree_test

import (
	merkletree "github.com/reactivejson/merkleTree/internal/merkle"
	"github.com/reactivejson/merkleTree/internal/merkle/hash"
	"math/rand"
	"testing"
)

func benchmarkMerkleTree(n int, b *testing.B) {
	// create a slice of n byte arrays with each element having a size of 32 bytes
	data := make([][]byte, n)
	for i := 0; i < n; i++ {
		data[i] = make([]byte, 32)
	}
	// create a new MerkleTree instance
	tree, err := merkletree.NewTree(data, hash.NewBlake3())
	if err != nil {
		b.Fatal(err)
	}
	// run the loop b.N times
	for i := 0; i < b.N; i++ {
		// generate the proof for a random data element
		index := rand.Intn(n)
		proof, err := tree.GenerateMProof(data[index])
		if err != nil {
			b.Fatal(err)
		}
		verified, err := merkletree.VerifyMProof(data[index], proof, tree.MerkleRoot(), blake3)
		if err != nil {
			b.Fatal(err)
		}
		_ = verified
	}
}

func BenchmarkMerkleTree10(b *testing.B)      { benchmarkMerkleTree(10, b) }
func BenchmarkMerkleTree100(b *testing.B)     { benchmarkMerkleTree(100, b) }
func BenchmarkMerkleTree1000(b *testing.B)    { benchmarkMerkleTree(1000, b) }
func BenchmarkMerkleTree10000(b *testing.B)   { benchmarkMerkleTree(10000, b) }
func BenchmarkMerkleTree100000(b *testing.B)  { benchmarkMerkleTree(100000, b) }
func BenchmarkMerkleTree1000000(b *testing.B) { benchmarkMerkleTree(1000000, b) }
