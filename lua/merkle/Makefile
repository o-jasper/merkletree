
default: test

test: test_sha2 test_merkle test_statement

# Test the sha2 stuff.
test_sha2: test_sha2_came_with test_sha2_more test_sha2_bin224 test_sha2_bin256

test_sha2_came_with:   # Tests it came with.
	lua test/sha2.lua

test_sha2_more:
	lua test/sha2_more.lua

test_sha2_bin256:  # Compares lua to program sha256sum
	sh test/test.sh 256 | head -n 10

test_sha2_bin224:  # idem
	sh test/test.sh 224 | head -n 10

# TODO compare luapower sha2? (doesnt work for me right now,
# luapower-all/bin/linux64/?.so is in `package.cpath`

# Test the merkle tree
test_merkle: test_merkle_devrandom

test_merkle_devrandom:
	lua Merkle/test/merkle.lua /dev/random

test_merkle_lua:
	lua Merkle/test/merkle.lua

# TODO test merkle version stuff of statement, identical-input.
test_statement:  # Just runs it, really.
	lua statement/test/statement.lua

# test_merkle_bin  ... TODO combine with go test?
