
default: test build

clean:
	rm test_merkletree path_chunk_n_root merkletree.so

test: test_merkle test_signed test_pubkey test_signed_negative test_trie

test_merkle:
	GOPATH=`pwd` go run src/test_merkletree.go

test_signed:
	GOPATH=`pwd` go run src/test_signed_merkletree.go

test_signed_negative:
	GOPATH=`pwd` go run src/test_signed_merkletree.go -negative true

test_pubkey:
	GOPATH=`pwd` go run src/test_pubkey.go

test_trie:
	GOPATH=`pwd` go run src/test_trie.go

test_trie_merkle:
	echo TODO
# GOPATH=`pwd` go run src/test_trie_easy.go

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
