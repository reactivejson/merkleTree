package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	merkletree "github.com/reactivejson/merkleTree/internal/merkle"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

type TreeRequest struct {
	Data []string `json:"data"`
	Name string   `json:"name"`
}
type ProofRequest struct {
	Data string `json:"data"`
	Name string `json:"name"`
}
type UpdateLeafReq struct {
	Data  string `json:"data"`
	Name  string `json:"name"`
	Index uint64 `json:"index"`
}

var trees = make(map[string]*merkletree.MerkleTree)

func byteArray(req TreeRequest) [][]byte {
	var data [][]byte
	for _, it := range req.Data {
		data = append(data, []byte(it))
	}
	return data
}

var hashing = merkletree.NewBlake3()

// @Summary Create a new Merkle tree
// @Description Creates a new Merkle tree with the given data
// @Tags Merkle trees
// @Accept  json
// @Produce  json
// @Param tree body TreeRequest true "The data and name for the new Merkle tree"
// @Success 200 {string} string	""
// @Failure 400 {object} ErrorResponse
// @Router /create [post]
func CreateTree(c *gin.Context) {
	var data TreeRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the tree
	tree, err := merkletree.NewTree(byteArray(data), hashing)

	if err != nil {
		c.Error(err)
	} else {
		trees[data.Name] = tree
		c.JSON(http.StatusOK, "")
	}
}

// @Summary Verify a Merkle proof
// @Description Verifies a Merkle proof for a given data value
// @Tags Merkle trees
// @Accept  json
// @Produce  json
// @Param proof body ProofRequest true "The data, name, and proof to verify"
// @Success 200 {object} VerificationResponse
// @Failure 400 {object} ErrorResponse
// @Router /verify [post]
func VerifyProof(c *gin.Context) {
	var data ProofRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the tree
	tree, ok := trees[data.Name]

	if !ok {
		c.Error(fmt.Errorf("no tree found %v", data.Name))
	} else {
		_, verified, _ := verify(tree, []byte(data.Data))
		c.JSON(http.StatusOK, gin.H{"verified": verified})
	}
}

func VisualizeProof(c *gin.Context) {
	var data ProofRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the tree
	tree, ok := trees[data.Name]

	if !ok {
		c.Error(fmt.Errorf("no tree found  %v", data.Name))
	} else {
		// Generate a proof for data
		proof, err := tree.GenerateProof([]byte(data.Data))
		if err != nil {
			c.Error(fmt.Errorf("failed to generate proof  %v", data.Name))
		}
		graph := tree.VisualProof(proof, new(merkletree.StringFormatter), nil)
		writeVisual(data.Name, graph)
		c.JSON(http.StatusOK, gin.H{"graph": graph})
	}
}

// @Summary Update a Merkle tree leaf
// @Description Updates a leaf node in a Merkle tree with new data
// @Tags Merkle trees
// @Accept  json
// @Produce  json
// @Param update body UpdateLeafReq true "The name, index, and new data for the leaf to update"
// @Success 200 {string} string	""
// @Failure 400 {object} ErrorResponse
// @Router /update [put]
func UpdateLeaf(c *gin.Context) {
	var data UpdateLeafReq
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the tree
	tree, ok := trees[data.Name]

	if !ok {
		c.Error(fmt.Errorf("no tree found  %v", data.Name))
	} else {
		tree.UpdateLeaf(data.Index, []byte(data.Data))
		c.JSON(http.StatusOK, "")
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// VerificationResponse represents a response to a Merkle proof verification request
type VerificationResponse struct {
	Verified bool `json:"verified"`
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		log.Printf("Error %v", err)
	}

	c.JSON(http.StatusBadRequest, "")
}

func verify(tree *merkletree.MerkleTree, data []byte) (*merkletree.MerkleProof, bool, error) {
	// Fetch the root hash of the tree
	root := tree.Root()

	// Generate a proof for data
	proof, err := tree.GenerateProof(data)
	if err != nil {
		return nil, false, err
	}

	// Verify the proof for 'Baz'
	verified, err := merkletree.VerifyMProofUsing(data, proof, root, hashing)

	return proof, verified, err
}

func writeVisual(f, s string) {
	absPath, err := filepath.Abs("./visual")
	if err != nil {
		fmt.Println("failed to writeVisual file path: ", err)
	}
	err = os.WriteFile(absPath+"/"+f+".dot", []byte(s), 0644)
	if err != nil {
		fmt.Println("failed to writeVisual file: ", err)
	}
}
