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

### Project layout

This layout is following pattern:

```text
merkleTree
└───
    ├── .github
    │   └── workflows
    │     └── go.yml
    ├── cmd
    │   └── main.go
    ├── internal
    │   └── merkleTree.go
    ├── internal
    │   └── merkleProof.go
    ├── internal
    │   └── hash
    │     └── blake3.go
    ├── build
    │   └── Dockerfile
    ├── Makefile
    ├── README.md
    └── <source packages>
```

## Setup

### Getting started
merkle-tree is available in github
[merkle-tree](https://github.com/reactivejson/merkleTree)

```shell
go get github.com/reactivejson/merkleTree
```

#### Run
```shell
go run cmd/main.go
```

#### Build
```shell
make build
```
#### Testing
```shell
make test
```
### Build docker image:

```bash
make docker-build
```
This will build this application docker image so-called merkle-tree

## Benchmarking

We have run benchmarks using the `go test` command with the following options: `-bench=Bench -benchtime=5s -benchmem`.

### Results
goos: windows
goarch: amd64

cpu: 11th Gen Intel(R) Core(TM) i7-1185G7 @ 3.00GHz

| Function                    | Iterations   | Execution Time   | Memory Allocation   | Allocations Per Operation |
|-----------------------------|--------------|------------------|---------------------|---------------------------|
| BenchmarkMerkleTree10-8     | 4907882      | 1213 ns/op       | 760 B/op            | 16 allocs/op              |
| BenchmarkMerkleTree100-8    | 3520026      | 1769 ns/op       | 1272 B/op           | 25 allocs/op              |
| BenchmarkMerkleTree1000-8   | 2718508      | 2374 ns/op       | 1768 B/op           | 34 allocs/op              |
| BenchmarkMerkleTree10000-8  | 1714077      | 3344 ns/op       | 2458 B/op           | 46 allocs/op              |
| BenchmarkMerkleTree100000-8 | 1598781      | 3777 ns/op       | 2975 B/op           | 55 allocs/op              |
| MerkleTreeOneMillion        | 1111663      | 4552 ns/op       | 3731 B/op           | 69 allocs/op              |


The `ns/op` column shows the average number of nanoseconds per operation. The `B/op` column shows the average number of bytes allocated per operation. The `allocs/op` column shows the average number of memory allocations per operation.

The execution time is measured in nanoseconds per operation, and it indicates how long it takes to execute each operation. The execution time is calculated by multiplying the number of iterations by the average time per iteration.

Note that the execution time may vary depending on the hardware and software configuration of the system running the benchmarks.

## Merkle Tree Package
This is a Go package that provides a Merkle tree data structure implementation.

### Usage
The package provides the merkletree package that contains the following functions:

#### NewTree(data [][]byte, hash HashType) (*MerkleTree, error)
This function creates a new MerkleTree struct that represents a Merkle tree of the given data using the specified HashType.

#### GenerateMProof(data []byte) (*MerkleProof, error)
This function generates a Merkle proof for a given data element. It returns a MerkleProof struct.

#### MerkleRoot() []byte
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
