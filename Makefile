
default: test build

clean:
	rm test_merkletree path_chunk_n_root merkletree.so

test:
	GOPATH=`pwd` go run src/test_merkletree.go

data:
	GOPATH=`pwd` go run src/path_chunk_n_root.go

build: test_merkletree path_chunk_n_root

test_merkletree: src/test_merkletree.go src/merkletree/merkletree.go src/common/common.go
	GOPATH=`pwd` go build src/test_merkletree.go

path_chunk_n_root: src/path_chunk_n_root.go src/merkletree/merkletree.go src/common/common.go
	GOPATH=`pwd` go build src/path_chunk_n_root.go

buildso: merkletree.so

# No idea on your milage.. Probably needs more work.
merkletree.so: src/merkletree/merkletree.go
	gccgo src/merkletree/merkletree.go -c -o merkletree.so
