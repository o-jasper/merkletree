See the `Merkle/` directory for the actual merkle tree api.
The `.H`, `.H2` members allow for selecting a hash function at
creation/"class deriving".

`sha2` is the sha2 by Roberto Ierusalimschy
[as on lua-users.org](http://lua-users.org/wiki/SecureHashAlgorithm)
([archive](https://archive.is/sJrRo)). It hash `.hash224` and `.hash256`.

**TODO**, it has `.new256` too, probably will want to improve on that.

`make test` runs the tests.
