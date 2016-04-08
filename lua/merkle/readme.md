### Sha2
`sha2.lua` is the sha2 originally by Roberto Ierusalimschy
[as on lua-users.org](http://lua-users.org/wiki/SecureHashAlgorithm)
([archive](https://archive.is/sJrRo)), and mostly only organizationally
altered by me. It has `.sha224(data)` and `.sha256(data)` to use plainly
as with data.

It also has `.Sha256` and `.Sha224` are the class versions, where
`:add(data)` can be used to add data, and `:close()` finalizes and
returns

### Merkle tree
See the `Merkle/` directory for the merkle tree api.
The `.H`, *and* `.H2` members allow for selecting a hash function
(like the above) at creation/"class deriving".

### `statement/`
Statements, roughly you can put in trees, and it produces a "root hash",
optionally with a nonce. The underlying system it uses can be
straight-from-hashes or with merkle tree. How to go depends on what result
you want.

Again, see the directory itself for specifics. (TODO)

## Testing
`make test` runs the tests. Sha2 has example tests, and tests
comparing it with the commands `sha224sum` and `sha256sum`.

Merkle trees are merely tested by making the tree and true/false proofs,
and checking for no false negatives/positives.

Statements are tested, no testing for false positives yet.

It is important to know what is actually tested. These tests run on
random data and examples. It might be good to look for edge cases.

## TODO
* `sha2.lua` needs `bit32`, which appears not available on `luajit`.

* When luajit, compare that too.

* The Merkle statements are derived from the merkle tree, thus it
  can be tested whether the merkle aspect infacts works properly.

* Test against statements false positives.

* Compare with the go implementation.
