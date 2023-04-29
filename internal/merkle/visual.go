package merkletree

import (
	"fmt"
	"math"
	"strings"
)

// Visual creates a Visual representation of the tree.  It is generally used for external presentation.
// This takes two optional formatters for []byte input: the first for leaf input and the second for branches.
func (t *MerkleTree) Visual(lf Formatter, bf Formatter) string {
	return t.visual(nil, nil, nil, lf, bf)
}

// VisualProof creates a Visual representation of the tree with highlights for a proof.  It is generally used for external presentation.
// This takes two optional formatters for []byte input: the first for leaf input and the second for branches.
func (t *MerkleTree) VisualProof(proof *MerkleProof, lf Formatter, bf Formatter) string {
	if proof == nil {
		return t.Visual(lf, bf)
	}

	// Find out which nodes are used in our proof
	valueIndices := make(map[uint64]int)
	proofIndices := make(map[uint64]int)
	rootIndices := make(map[uint64]int)

	if proof != nil {
		index := proof.Index + uint64(math.Ceil(float64(len(t.nodes))/2))
		valueIndices[proof.Index] = 1

		for range proof.Hashes {
			proofIndices[index^1] = 1
			index /= 2
		}

		numRootNodes := uint64(math.Exp2(math.Ceil(math.Log2(float64(len(t.data))))-float64(len(proof.Hashes))+1)) - 1
		for i := uint64(1); i <= numRootNodes; i++ {
			rootIndices[i] = 1
		}
	}

	return t.visual(rootIndices, valueIndices, proofIndices, lf, bf)
}
func (t *MerkleTree) visual(rootIndices, valueIndices, proofIndices map[uint64]int, lf, bf Formatter) string {
	if lf == nil {
		lf = new(TruncatedHexFormatter)
	}
	if bf == nil {
		bf = new(TruncatedHexFormatter)
	}

	var builder strings.Builder
	builder.WriteString("digraph MerkleTree {")
	builder.WriteString("rankdir = TB;")
	builder.WriteString("node [shape=rectangle margin=\"0.2,0.2\"];")
	empty := make([]byte, len(t.nodes[1]))
	dataLen := len(t.data)
	valuesOffset := int(math.Ceil(float64(len(t.nodes)) / 2))
	var nodeBuilder strings.Builder
	nodeBuilder.WriteString("{rank=same")
	for i := 0; i < valuesOffset; i++ {
		if i < dataLen {
			// Value
			builder.WriteString(fmt.Sprintf("\"%s\" [shape=oval", lf.Format(t.data[i])))
			if valueIndices[uint64(i)] > 0 {
				builder.WriteString(" style=filled fillcolor=\"#00FFFF\"")
			}
			builder.WriteString("];")

			builder.WriteString(fmt.Sprintf("\"%s\"->%d;", lf.Format(t.data[i]), valuesOffset+i))

			nodeBuilder.WriteString(fmt.Sprintf(";%d", valuesOffset+i))
			builder.WriteString(fmt.Sprintf("%d [label=\"%s\"", valuesOffset+i, bf.Format(t.nodes[valuesOffset+i])))
			if proofIndices[uint64(i+valuesOffset)] > 0 {
				builder.WriteString(" style=filled fillcolor=\"#FFFF00\"")
			} else if rootIndices[uint64(i+valuesOffset)] > 0 {
				builder.WriteString(" style=filled fillcolor=\"#C0C0C0\"")
			}
			builder.WriteString("];")
			if i > 0 {
				builder.WriteString(fmt.Sprintf("%d->%d [style=invisible arrowhead=none];", valuesOffset+i-1, valuesOffset+i))
			}
		} else {
			// Empty leaf
			builder.WriteString(fmt.Sprintf("%d [label=\"%s\"", valuesOffset+i, bf.Format(empty)))
			if proofIndices[uint64(i+valuesOffset)] > 0 {
				builder.WriteString(" style=filled fillcolor=\"#FFFF00\"")
			} else if rootIndices[uint64(i+valuesOffset)] > 0 {
				builder.WriteString(" style=filled fillcolor=\"#C0C0C0\"")
			}
			builder.WriteString("];")
			builder.WriteString(fmt.Sprintf("%d->%d [style=invisible arrowhead=none];", valuesOffset+i-1, valuesOffset+i))
			nodeBuilder.WriteString(fmt.Sprintf(";%d", valuesOffset+i))
		}
		if dataLen > 1 {
			builder.WriteString(fmt.Sprintf("%d->%d;", valuesOffset+i, (valuesOffset+i)/2))
		}
	}
	nodeBuilder.WriteString("};")
	builder.WriteString(nodeBuilder.String())

	// Add branches
	for i := valuesOffset - 1; i > 0; i-- {
		builder.WriteString(fmt.Sprintf("%d [label=\"%s\"", i, bf.Format(t.nodes[i])))
		if rootIndices[uint64(i)] > 0 {
			builder.WriteString(" style=filled fillcolor=\"#C0C0C0\"")
		} else if proofIndices[uint64(i)] > 0 {
			builder.WriteString(" style=filled fillcolor=\"#FFFF00\"")
		}
		builder.WriteString("];")
		if i > 1 {
			builder.WriteString(fmt.Sprintf("%d->%d;", i, i/2))
		}
	}
	builder.WriteString("}")
	return builder.String()
}

// Formatter formats a []byte in to a string.
// It is used by Visual() to provide users with the required format for the graphical display of their Merkle trees.
type Formatter interface {
	// Format
	Format([]byte) string
}

// TruncatedHexFormatter shows only the first and last two bytes of the value.
type TruncatedHexFormatter struct{}

// Format formats a value as truncated hex, showing the first and last four characers of the hex string.
func (f *TruncatedHexFormatter) Format(data []byte) string {
	return fmt.Sprintf("%4x…%4x", data[0:2], data[len(data)-2:])
}

// HexFormatter shows the entire value.
type HexFormatter struct{}

// Format formats a value as a full hex string.
func (f *HexFormatter) Format(data []byte) string {
	return fmt.Sprintf("%0x", data)
}

// StringFormatter shows the entire value as a string.
type StringFormatter struct{}

// Format formats a value as a UTF-8 string.
func (f *StringFormatter) Format(data []byte) string {
	return string(data)
}
