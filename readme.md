## Merkle tree creator
The merkle tree should eat 'chunks' and put them on the leaves. It is able to
mark chunks(leafs) as interesting, so it doesnt have to actually store the whole thing.

This is useful for where large portions of the merkle tree are infact simply
calculated entries.

The Ethereum contract merely has to check if the leaf, path and root connect
correctly.

#### Uses

* Lots-a-stuff, i dont know well about, torrents use them for instance.

* Allowing Ethereum to have access to data. Also 'pre-emptively' for data that
  you dont think Ethereum contracts have any use for.

* [Hanging blocks](http://o-jasper.github.io/blog/2014/06/03/hanging_blocks.html)
  of various sorts.

* [Dropbox example](https://github.com/jorisbontje/cll-sim/blob/master/examples/decentralized-dropbox.cll), i suppose.

### Cumulative merkle-tree generation functions
These are functions of the Merkle tree, using stuff in `src/merkle/`
and `src/hash_extra/`, not any of the other stuff. The trie doesnt work.

not the signed merkle tree

    func NewMerkleTreeGen(hasher Hasher, include_index bool) *MerkleTreeGen
    
Creates an object that gathers chunks, creating the Merkle tree on the way.

`hasher.H(leaf)` must produce a hash, `hasher.HwI(i,leaf)` a hash-with-an-index
and `hasher.H_U2(a,b)` produces a parent hash from two.

`include_index` determines if it is just a set of objects, of if they're indexed.

    func (gen *MerkleTreeGen) AddChunk(leaf []byte, interest bool) *MerkleNode
    
Allows `MerkleTreeGen` to do its thing, adding a `leaf` of data. 
It returns `*MerkleNode`, which can be used to create those paths, *if*
`interest == true`.

    func (gen *MerkleTreeGen) Finish() *MerkleNode

After calling this you can use the returned `*MerkleNode` as if you are
finished, it can be used to get at the root hash (`.Hash`). You can continue, 
however, but the paths made from the node then go past that hash.

    func (node *MerkleNode) Path() [][sha256.Size]byte
    func (node *MerkleNode) ByteProof() [][sha256.Size]byte
    
Makes a path from a merkle node to the top, so that it can be proven that the
checksum of the a leaf corresponds to the root checksum.

    func (node* MerkleNode) VerifyH(hasher Hasher, Hroot, Hleaf HashResult) int8
    func (node* MerkleNode) Verify(hasher Hasher, Hroot HashResult, leaf []byte) int8
    func (node* MerkleNode) VerifyWithIndex(hasher Hasher, Hroot HashResult, i uint64, leaf []byte) int8

Just the on-build, **not** the typically used verifier, which is in the `Hasher` in
the section below.

Returns one of `Correct`, or ways it is wrong; 	`WrongDataPath`, `WrongDataLeaf`,
`WrongDataRoot`, `WrongSigPath`, `WrongSigLeaf`, `WrongSigRoot`, `WrongSomeThing`,
`WrongSig`.

### Hasher functions
A hasher takes a hash function, and produces `.H()`, `.HwI`, for hashing, and
`.H_2`, and `.H_U2` to combine two hashes.

One way to create one is `hash_extra.Hasher{sha256.New()}` i.e. just give it the
less-featured hash to use.

Additionally, hashers can also verify proofs;

    func (hasher Hasher) MerkleVerifyH(root, Hleaf HashResult, path []HashResult) bool
    func (hasher Hasher) MerkleVerify(root HashResult, leaf []byte, path []HashResult) bool
    func (hasher Hasher) MerkleVerifyWithIndex(root HashResult, i uint64, leaf []byte, path []HashResult) bool

Return whether verification succeeded given the information available.

    func (hasher Hasher) MerkleExpectedRoot(H_leaf HashResult, path []HashResult) HashResult

Returns what the root is *expected* to be.

## TODO

* Improve the Hasher; instead of `H_2` or `H_U2` appending,
  use `bitwise_xor(a, bitwise_not(b))` or something, it is faster.
  (matters particularl for Ethereum entities)

* Make a corresponding contract that merely serves a verifying function given a
  root, leaf checksum and path. (NOTE: it used to work?)
 
  + Then do the dropbox example.
  + Possibly namereg example.

* The above docs could look better.
* Negative result tests.

* What the hell is `src/bin/path_chunk_n_root.go` for?
