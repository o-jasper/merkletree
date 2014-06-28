
default: test build

test:
	GOPATH=`pwd` go run src/test_merkletree.go

clean:
	rm test_merkletree merkletree.so

build: test_merkletree

test_merkletree: src/test_merkletree.go src/merkletree/merkletree.go
	GOPATH=`pwd` go build src/test_merkletree.go

buildso: merkletree.so

# No idea on your milage.. Probably needs more work.
merkletree.so: src/merkletree/merkletree.go
	gccgo src/merkletree/merkletree.go -c -o merkletree.so
