## `Merkle.Tree`

Subelements `.H(data)` must produce the hash of data, `.H2(a, b)` produces the
hash of a pair. (which combines pairs for the tree) These can be added to the object
or you can derive from this thing and add them in the metatable.

* `Tree:add(data, keep_proof)` adds `data` as a leaf, `keep_proof` indicates whether
  the proof should be kept, in which case it is returned as `Node` which can
  construct proofs

  Similarly `:add_H(H, keep_proof)`, but there you do the hash for it.

* `Tree:add_key(keydata, data, keep_proof)` is like `:add(keydata .. data, keep_proof)`,
  instead that *if* `keep_proof`, then `.kept_keys[keydata]` refers to the Merkle
  node, so that the proofs for a key can readily be acquired.

*  `Tree:close()` finishes it and returns the root.

  You *can* continue and close again, changing the root, and creating a lobsided
  tree.

## `Merkle.Node`
These are the handle to proof-creators, they're the objects returned
by `Tree:add`, `..:add_H`, `:add_key`.

*  `Node:produce_proof()` returns the proof, should work on the root hash *at that time*.

  Note: can supply a list to start with, `:produce_proof(ret_list)`, and it'll
  add the result to the list.

## `Merkle.Verify`

* `Verify:verify(root, proof, leaf)` verifies the proof given the `root` hash and `leaf` data.

* `Verify:verify_H(root, proof, leaf_H)` is a version requiring you do to the hash.

* `Verify:verify_key(root, proof, key, leaf)` keyed version.

* `Verify:expect_root(proof, leaf)` returns the expected root for it to be correct.
  (also a `:expect_root_H`)
