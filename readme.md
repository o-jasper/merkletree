## Merkle tree creator
The merkle tree should eat 'chunks' and put them on the leaves. It is intended
to be able to mark chunks as interesting, so it doesnt have to actually store
the whole thing.

This is useful for where large portions of the merkle tree are infact simply
calculated entries.

The Ethereum contract merely has to check.

#### Uses

* Lots-a-stuff, i dont know well about, torrents use them for instance.

* Allowing Ethereum to have access to data. Also 'pre-emptively' for data that
  you dont think Ethereum contracts have any use for.

* [Hanging blocks](http://o-jasper.github.io/blog/2014/06/03/hanging_blocks.html)
  of various sorts.

* [Dropbox example](https://github.com/jorisbontje/cll-sim/blob/master/examples/decentralized-dropbox.cll), i suppose.

#### Additional
An idea is to take probabilities for chunks and make a somewhat lobsided tree
that minimizes the average length of proving merkle paths. The checking mechanism
doesnt care about the shape of the merkle tree, so it can be added later,
also, `merkletree.Finish` can sort-of be used prematurely for lobsidedness.

On the other hand, off-chain, the efficiency in memory is better if things you want
to keep track of are bunched together.

### Functions

    func NewMerkleTreeGen() *MerkleTreeGen
    
Creates an object that gathers chunks, creating the Merkle tree on the way.

    func (gen *MerkleTreeGen) AddChunk(chunk []byte, interest bool) *MerkleNode
    
Allows `MerkleTreeGen` to do its thing, adding a `chunk` of data. 
It returns `*MerkleNode`, which can be used to create those paths, *if*
`interest == true`.

    func (gen *MerkleTreeGen) Finish() *MerkleNode

After calling this you can use the returned `*MerkleNode` as if you are
finished, it can be used to get at the root hash (`.Hash`). You can continue, 
however, but the paths mad from the node then go past that hash.

    func (node *MerkleNode) Path() [][sha256.Size]byte
    
Makes a path from a merkle node to the top, so that it can be proven that the
checksum of the a leaf corresponds to the root checksum.

    func ExpectedRoot(H_leaf [sha256.Size]byte, path [][sha256.Size]byte) [sha256.Size]byte

Returns the root expected, based on the leaf hash, and path.

    func CorrectRoot(root [sha256.Size]byte, leaf []byte, path [][sha256.Size]byte) bool
    
Returns whether the root is correct, given the leaf chunk and path.

**The following two** are sort-of alternative ways to use paths, they use parts
of the constructed tree. However, the above are provided because you need a way
to get the data is a simp;e binary format.

    func (node *MerkleNode) IsValid(recurse int32) bool

Tells you if the known tree upward from the given merkle node by the given
recursions are valid. `recurse < 0` means that it will recurse all the way.

    func (node *MerkleNode) CorrespondsToChunk(chunk []byte) bool

Tells you that the `*MerkleNode` is 1) a leaf, and 2) corresponds to the chunk.

**Some additional functions** are `H`, `H_2`, which are the how `sha256.Sum256`
is modified to have the additional right/left and uninteresting/interesting 
information.

## TODO

* Make a corresponding contract that merely serves a verifying function given a
  root, leaf checksum and path.
 
  + Then do the dropbox example.
  + Possibly namereg example.
