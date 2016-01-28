
default: test

test: test_sha2 test_merkle

# Test the sha2 stuff.
test_sha2: test_sha2_came_with test_sha2_bin

test_sha2_came_with:   # Tests it came with.
	lua test/sha2.lua

test_sha2_bin:  # Compares to the function. (NOTE: only the 256 i think tho..
	sh test/test.sh | head -n 100

# Test the merkle tree
test_merkle: test_merkle_lua

test_merkle_lua:
	lua Merkle/test/merkle.lua

# test_merkle_bin  ... TODO combine with go test?
