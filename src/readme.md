## hash_extra/
Adds set of functions useful for merkle trees, and can also check proofs.

## merkle/
Main merkle tree lib. Needs `hash_extra/`

(for those extra features, like `H_2`, combining two hashes)

## signed_merkle/
The idea is/was:

* Send data.
* Send nonce, demand root of all chunks signed with nonce.
* Demand proof of some particular signed chunks. And the signed chunk  itself.

So that you have to have store all the data *and* the private key, because you don't know
what chunk will be asked to be proven, and otherwise you need to get all data send past you.

Not sure on wisdom of that idea..

## trie_easy/
Don't use it.. Forgot about it myself..

## test/
Test of the different things here. `test/common/` contains functions,
like random generating chunks aiding the tests.

## bin/path_chunk_n_root.go
Probably useless.. Keeping it around to run into something useful.
