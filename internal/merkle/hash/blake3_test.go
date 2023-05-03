package hash_test

import (
	"encoding/hex"
	"fmt"
	hash2 "github.com/reactivejson/merkleTree/internal/merkle/hash"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

func stringToByte(input string) []byte {
	x, err := hex.DecodeString(strings.TrimPrefix(input, "0x"))
	if err != nil {
		panic(err)
	}
	return x
}
func TestBlake3(t *testing.T) {
	tests := []struct {
		input []byte
		hash  []byte
	}{
		{
			input: []byte("Consensys"),
			hash:  stringToByte("37d279155d7afba864451532eb236103d43b8d410806322ea36be2b8f7731dfd"),
		},
	}

	hash := hash2.NewBlake3()
	for i, test := range tests {
		res := hash.Hash(test.input)
		fmt.Println(hex.EncodeToString(res))
		assert.Equal(t, test.hash, res, fmt.Sprintf("failed at test %d", i))
	}
}

func TestConcatHash(t *testing.T) {
	tests := []struct {
		input1 []byte
		input2 []byte
		input3 []byte
		hash   []byte
	}{
		{ // 0
			input1: []byte("Merle-tree"),
			input2: []byte("Blake3"),
			input3: []byte("Consensys"),
			hash: []byte{0xfd, 0xb6, 0x8f, 0x8b, 0x88, 0x59, 0xb0, 0xab, 0x23, 0x9a, 0xd5, 0x86, 0x6, 0xe1, 0xd3, 0x16,
				0xd0, 0x7, 0xa8, 0x3e, 0x86, 0x7d, 0xc0, 0x84, 0x2a, 0xf4, 0xcd, 0x99, 0x62, 0x14, 0x76, 0x51},
		},
	}

	hash := hash2.NewBlake3()
	for i, test := range tests {
		res := hash.Hash(test.input1, test.input2, test.input3)
		assert.Equal(t, test.hash, res, fmt.Sprintf("failed at test %d", i))
	}
}
