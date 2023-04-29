## Merkle Tree API
This is Go implementation of Merkle Tree data structure, 
which is a hash-based data structure that is useful for verifying the integrity of data.
It is used in many different contexts, including in blockchain technology.

We use Blake3 hashing algorithm.

## Endpoints
The API has a single endpoint:

### POST /create
Create a new Merkle tree
This endpoint accepts a JSON payload containing a list of data and name for the new Merkle tree.

Request Payload

Example request payload:
````json
{
  "data": ["Foo", "Bar", "Baz"],
  "name": "tree1"
}
````

### POST /verify
Verify a Merkle proof for a given data item
This endpoint accepts a JSON payload containing the data to be verified and name for the new Merkle tree.

Request Payload

````json
{
  "data": "Baz",
  "name": "tree1"
}

````


Response Payload
The response payload is a JSON object with a single field verified that is a boolean

Example response payload:

````json
{
  "verified": true
}
````

### PUT /update
Update a leaf node in a Merkle tree
This endpoint accepts a JSON payload containing The name, index, and new data for the leaf to update.

Request Payload

````json
{
  "data": "kk",
  "name": "tree1",
  "index": 2
}

````


## Benchmarking

We have run benchmarks using the `go test` command with the following options: `-bench=Bench -benchtime=5s -benchmem`.

### Results

| Function | Iterations | Execution Time | Memory Allocation | Number of Allocations |
| -------- | ---------- | -------------- | ----------------- | --------------------- |
| MerkleTree10 | 2354869 | 4105 ns/op | 760 B/op | 16 allocs/op |
| MerkleTree100 | 976814 | 5134 ns/op | 1272 B/op | 25 allocs/op |
| MerkleTree1000 | 818107 | 6514 ns/op | 1768 B/op | 34 allocs/op |
| MerkleTree10000 | 561963 | 10252 ns/op | 2463 B/op | 46 allocs/op |
| MerkleTree100000 | 470954 | 12149 ns/op | 3031 B/op | 56 allocs/op |

The `ns/op` column shows the average number of nanoseconds per operation. The `B/op` column shows the average number of bytes allocated per operation. The `allocs/op` column shows the average number of memory allocations per operation.

The execution time is measured in nanoseconds per operation, and it indicates how long it takes to execute each operation. The execution time is calculated by multiplying the number of iterations by the average time per iteration. Here is the execution time for each function:

| Function | Execution Time |
| -------- | -------------- |
| MerkleTree10 | 9.67 s |
| MerkleTree100 | 5.01 s |
| MerkleTree1000 | 5.33 s |
| MerkleTree10000 | 5.76 s |
| MerkleTree100000 | 5.72 s |

Note that the execution time may vary depending on the hardware and software configuration of the system running the benchmarks.

## Merkle Tree Package
This is a Go package that provides a Merkle tree data structure implementation.

### Usage
The package provides the merkletree package that contains the following functions:

#### NewTree(data [][]byte, hash HashType) (*MerkleTree, error)
This function creates a new MerkleTree struct that represents a Merkle tree of the given data using the specified HashType.

#### GenerateProof(data []byte) (*MerkleProof, error)
This function generates a Merkle proof for a given data element. It returns a MerkleProof struct.

#### Root() []byte
This function returns the Merkle root hash.

#### UpdateLeaf(index uint64, newData []byte) error
This function updates the leaf at the given index with new data. It returns an error if the index is out of bounds.

#### VerifyMProof(data []byte, proof *MerkleProof, root []byte) (bool, error)
This function verifies a given Merkle proof against a Merkle root hash using the Blake3 hashing algorithm. It returns a boolean value indicating whether the proof is valid or not.

### Types
The package provides the following types:

#### MerkleTree struct
This struct represents a Merkle tree data structure. It contains the following fields:

hash HashType: A HashType object that represents the hashing algorithm used to generate the Merkle tree.
data [][]byte: A slice of byte slices that contains the original data elements used to generate the Merkle tree.
nodes [][]byte: A slice of byte slices that contains the nodes of the Merkle tree.
#### MerkleProof struct
This struct represents a Merkle proof. It contains the following fields:

Hashes [][]byte: A slice of byte slices that contains the hashes of the nodes on the proof path.
Index uint64: An integer that represents the index of the data element that the proof is for.

### Hashing
The package provides a HashType interface that defines the methods required for a hashing algorithm to be used with the MerkleTree struct.
The package includes a Blake3 hashing algorithm implementation which we used for this implementation.