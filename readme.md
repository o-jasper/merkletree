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

#### Additional
An idea is to take probabilities for leafs and make a somewhat lobsided tree
that minimizes the average length of proving merkle paths. The checking mechanism
doesnt care about the shape of the merkle tree, so it can be added later,
also, `merkletree.Finish` can sort-of be used prematurely for lobsidedness.

On the other hand, off-chain, the efficiency in memory is better if things you want
to keep track of are bunched together.

### Functions
These are functions of the Merkle tree, using stuff in `src/merkle/`
and `src/hash_extra/`, not any of the other stuff.

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

Returns one of `Correct`, or ways it is wrong; 	`WrongDataPath`, `WrongDataLeaf`,
`WrongDataRoot`, `WrongSigPath`, `WrongSigLeaf`, `WrongSigRoot`, `WrongSomeThing`,
`WrongSig`.


    func ExpectedRoot(H_leaf [sha256.Size]byte, path [][sha256.Size]byte) [sha256.Size]byte

Returns the root expected, based on the leaf hash, and path.

    func Verify(root [sha256.Size]byte, leaf []byte, path [][sha256.Size]byte) bool
    
Returns whether the root is correct, given the leaf chunk and path.
(`VerifyH` requires you to `H(leaf)`)

**The following three** are sort-of alternative ways to use paths, they use parts
of the constructed tree. However, the above are provided because you need a way
to get the data is a simple binary format.

    func (node* MerkleNode) Verify(Hroot [sha256.Size]byte, leaf []byte) bool

Verifies the whole thing, given a node that was indicated as interesting, the
root hash and the leaf hash. (`VerifyH` requires you to `H(leaf)`)

    func (node *MerkleNode) IsValid(recurse int32) MerkleNode, bool

Tells you if the known tree upward from the given merkle node by the given
recursions are valid. `recurse < 0` means that it will recurse all the way.

    func (node *MerkleNode) CorrespondsH(leaf HashResult) bool
    func (node *MerkleNode) Corresponds(leaf []byte) bool
    func (node *MerkleNode) CorrespondsWithIndex(hasher Hasher, i uint64, leaf []byte) bool

Tells you that the `*MerkleNode` is 1) a leaf, and 2) corresponds to the
hash or (positioned)chunk.

**Some additional functions** are `H`, `H_2`, which are the how `sha256.Sum256`
is modified to have the additional right/left and uninteresting/interesting 
information.

## TODO

* An alternative Hasher; instead of `H_2` or `H_U2` appending,
  use `bitwise_xor(a, bitwise_not(b))` or something, it is faster.
  (matters for Ethereum entity)

* Make a corresponding contract that merely serves a verifying function given a
  root, leaf checksum and path. (NOTE: it used to work?)
 
  + Then do the dropbox example.
  + Possibly namereg example.

* The above docs could look better.
* Negative result tests.
