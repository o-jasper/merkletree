
package merkle

import "hash"
import . "hash_extra"

// Note.. for instances of hash.Hash i am relying on copy-value..
// _seems_ like the spec is unclear on this? (not read enough)

// A node of Merkle tree, note that the below omits a lot.
type MerkleNode struct {
	Hash   hash.Hash
	Left   *MerkleNode
	Right  *MerkleNode
	Up     *MerkleNode  // Otherwise it is messy to create the path.
	Interest bool
}

func selective_new_MerkleNode(left, right *MerkleNode) *MerkleNode {
	n := MerkleNode{Hash : H_U2(left.Hash, right.Hash), Left : left, Right : right, Up:nil}
	n.Interest = left.Interest || right.Interest
	if n.Interest {
		left.Up, right.Up = &n, &n
	} else {
		n.Left, n.Right = nil, nil
	}
	return &n
}

func (node MerkleNode) effect_interest() {
	if !node.Interest { //No interest in branches.
		node.Left = nil
		node.Right = nil
	}
}

type MerkleTreePortion struct { //Subtree in there somewhere.
	Node   *MerkleNode
	Depth  int
}

//The algo is to make a 'mountain range'(forgot link)
// and put the latest together if we have them. Meanwhile, it creates a tree, but
// drops anything in which there is no interest.
type MerkleTreeGen struct {
	hash.Hash
	I uint64
	List []MerkleTreePortion
// Whether to prepend the index in the chunks.(otherwise order is not recorded at all)
	IncludeIndex bool
}

func NewMerkleTreeGen(h hash.Hash, include_index bool) MerkleTreeGen {
	return MerkleTreeGen{h, uint64(0), []MerkleTreePortion{}, include_index}
}

// Adds chunk where you calculated the hash, returning the leaf the current is on.
func (gen *MerkleTreeGen) AddChunkH(h hash.Hash, interest bool) *MerkleNode {
	gen.I += 1
	if len(gen.List) == 0 || gen.List[0].Depth != 1 {
		add_node := &MerkleNode{Hash:h, Left:nil, Right:nil, Up:nil, Interest:interest}

		list := []MerkleTreePortion{}
		list = append(list, MerkleTreePortion{Node:add_node, Depth:1})
		gen.List = append(list, gen.List...)
		return add_node
	} else {
		// assert gen.List[0].Depth == 1
		new_leaf := &MerkleNode{Hash:h, Left:nil, Right:nil, Interest:interest}

		new_node := selective_new_MerkleNode(new_leaf, gen.List[0].Node)  //Combine the two.
		gen.List[0] = MerkleTreePortion{Node:new_node, Depth:2}

		// Combine more, while equal depth.
		for len(gen.List) >= 2 && gen.List[1].Depth == gen.List[0].Depth {
			new_node := selective_new_MerkleNode(gen.List[0].Node, gen.List[1].Node)
			gen.List = gen.List[1:] // Cut off the first one.
			gen.List[0] = MerkleTreePortion{Node:new_node, Depth:gen.List[0].Depth + 1}
		}
		return new_leaf  //Return the leaf.
	}
}
// Calculates the hash for you.
func (gen *MerkleTreeGen) AddChunk(leaf []byte, interest bool) *MerkleNode {
	h := gen.Hash
	h.Reset()
	if gen.IncludeIndex {
		h.Write(getBytes(gen.I))
	}
	h.Write(leaf)
	return gen.AddChunkH(h, interest)
}

// Coerce the last parts together, returning the root.
// NOTE: you can 'finish' and then continue to make what you put in already
// becomes a bit of a 'lob' that takes longer Merkle paths.
func (gen *MerkleTreeGen) Finish() *MerkleNode {
	// assert len(gen.List) > 0
	for len(gen.List) >= 2  {
		new_node := selective_new_MerkleNode(gen.List[0].Node, gen.List[1].Node)
		gen.List = gen.List[1:]
		gen.List[0] = MerkleTreePortion{Node:new_node, Depth:gen.List[0].Depth}
	}
	return gen.List[0].Node
}

// Only checks internal consistency upward;
// Checks a merkle path, _except_ the root and leaf.
func (node *MerkleNode) IsValid(recurse int32) (*MerkleNode, bool) {
	switch {
	case node.Left != nil && node.Right != nil &&	H_U2(node.Left.Hash, node.Right.Hash) != node.Hash:
		return node, false
	case recurse == 0 || node.Up == nil:
		return node, true
	default:
		return node.Up.IsValid(recurse - 1)
	}
}

func (node *MerkleNode) CorrespondsToChunk(leaf []byte) bool {
	h := node.Hash
	h.Reset()
	h.Write(leaf)
	return node.CorrespondsToHash(h)
}

func (node *MerkleNode) CorrespondsToHash(h hash.Hash) bool {
	return h == node.Hash
}

// Used to run Data path errors into Sig (path) errors.
const Merkletree_NPathWrongs = int8(3)

const ( //NOTE: Not all functions will do all of them. Move some about signatures?
	Correct int8 = iota

	WrongDataPath
	WrongDataLeaf
	WrongDataRoot

	WrongSigPath
	WrongSigLeaf
	WrongSigRoot

	WrongSomeThing
	WrongSig
)

func (node* MerkleNode) VerifyH(Hroot, Hleaf hash.Hash) int8 {
	root, internal := node.IsValid(-1)
	switch {
	case !internal:                       return WrongDataPath
	case !node.CorrespondsToHash(Hleaf):  return WrongDataLeaf
	case !root.CorrespondsToHash(Hroot):  return WrongDataRoot
	default:                              return Correct
	}
}

func h(basis hash.Hash, data []byte) hash.Hash {
	h := basis
	h.Reset()
	h.Write(data)
	return h
}

func h_wi(basis hash.Hash, i uint64, data []byte) hash.Hash {
	h := basis
	h.Reset()
	h.Write(getBytes(i))
	h.Write(data)
	return h
}

func (node* MerkleNode) Verify(Hroot hash.Hash, leaf []byte) int8 {
	return node.VerifyH(Hroot, h(Hroot, leaf))
}

func (node* MerkleNode) VerifyWithIndex(Hroot hash.Hash, i uint64, leaf []byte) int8 {
	return node.VerifyH(Hroot, h_wi(Hroot, i, leaf))
}

// Calculated paths essentially make a compilation of the data needed to do the
// check. 
func (node *MerkleNode) ByteProof() [][]byte { //TODO []byte..
	ret, path := [][]byte{}, node.Path()
	for _, el := range path {	ret = append(ret, el.Sum([]byte{})) }
	return ret
}

func (node *MerkleNode) Path() []hash.Hash {
	switch {
	case node.Right != nil || node.Left != nil:  return nil
	case node.Up == nil:                       	 return nil
	default:                                     return node.Up.path(node)
	}
}

func (node *MerkleNode) path(from *MerkleNode) []hash.Hash {
	path := []hash.Hash{}
	if node.Up != nil {	path = node.Up.path(node) }
	if node.Right != from  && node.Right != from { 
		return nil
	}
	return append(path, node.Left.Hash)
}

//Calculate expected root, given the path.
func ExpectedRoot(H_leaf hash.Hash, path []hash.Hash) hash.Hash {
	x := H_leaf
	for i := range path {	x = H_U2(path[len(path) - i - 1], x) }
	return x
}

//Checks a root.
func VerifyH(root, Hleaf hash.Hash, path []hash.Hash) bool {
	return HashEqual(ExpectedRoot(Hleaf, path), root)
}
func Verify(root hash.Hash, leaf []byte, path []hash.Hash) bool {
	return VerifyH(root, h(root, leaf), path)
}
func VerifyWithIndex(root hash.Hash, i uint64, leaf []byte, path []hash.Hash) bool {
	return VerifyH(root, h_wi(root, i, leaf), path)
}
