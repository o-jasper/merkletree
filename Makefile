
default: test build

test:
	GOPATH=`pwd` go run src/test_merkletree.go

build:  # TODO .so or something?
	GOPATH=`pwd` go build src/test_merkletree.go
