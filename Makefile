
default: test build

clean:
	rm dist/test/erkletree path_chunk_n_root dist/merkletree.so

test: test_merkle test_signed test_pubkey test_signed_negative
#test_trie

test_merkle: src/test/merkletree.go
	echo ==== test_merkle;\
	GOPATH=`pwd` go run $<

test_signed: src/test/signed_merkletree.go
	echo ==== test_signed;\
	GOPATH=`pwd` go run $<

test_signed_negative: src/test/signed_merkletree.go
	echo ==== test_signed_negative;\
	GOPATH=`pwd` go run $< -negative true

test_pubkey: src/test/pubkey.go
	echo ==== test_pubkey;\
	GOPATH=`pwd` go run $<

test_trie: src/test/trie.go
	echo ==== test_trie
	GOPATH=`pwd` go run $<

test_trie_merkle:
	echo TODO
# GOPATH=`pwd` go run src/test_trie_easy.go

data:
	GOPATH=`pwd` go run src/path_chunk_n_root.go

build: dist/test/merkletree dist/path_chunk_n_root

# NOTE basically just there to check if differs form `go run ...`
dist/test/merkletree: src/test/merkletree.go src/merkle/merkletree.go src/merkle/merkle_common/common.go dist/test/
	GOPATH=`pwd` go build src/test/merkletree.go ; mv merkletree dist/test/

dist/test/: dist/
	mkdir dist/test

dist/:
	mkdir dist/

dist/path_chunk_n_root: src/bin/path_chunk_n_root.go src/merkle/merkletree.go src/merkle/merkle_common/common.go dist/
	GOPATH=`pwd` go build src/bin/path_chunk_n_root.go ; mv path_chunk_n_root dist/

buildso: merkletree.so

# No idea on your milage.. Probably needs more work.
dist/merkletree.so: src/merkle/merkletree.go dist/
	GOPATH=`pwd` gccgo src/merkle/merkletree.go -c -o dist/merkletree.so
