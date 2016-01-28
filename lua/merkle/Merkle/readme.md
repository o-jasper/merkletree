## Merkle `Tree`

Subelements `.H(data)` must produce the hash of data, `.H2(a, b)` produces the
hash of a pair. (which combines pairs for the tree) These can be added to the object
or you can derive from this thing and add them in the metatable.

* `Tree:add(data, keep_proof)` adds `data` as a leaf, `keep_proof` indicates whether
  the proof should be kept, in which case it is returned as `Node` which can
  construct proofs

  Similarly `:add_H(H, keep_proof)`, but there you do the hash for it.

*  `Tree:finish()` finishes it and returns the root.

  You *can* continue, changing the root, and creating a lobsided tree.

## Merkle `Node`
These are the handle to proof-creators, they're the objects returned by `Tree:add`, `..:add_H`.

*  `Node:produce_proof()` returns the proof, should work on the root hash *at that time*.

  Note: can supply a list to start with, `:produce_proof(ret_list)`, and it'll
  add the result to the list.

## `Verify`ing

* `Verify:verify(root, proof, leaf)` verifies the proof given the `root` hash and `leaf` data.

* `:verify_H(root, proof, leaf_H)` is a version requiring you do to the hash.

* `Verify:expect_root(proof, leaf)` returns the expected root for it to be correct.
  (also a `:expect_root_H`)
