
package merkletree

import (
	"crypto/sha256"
)

func FirstBit(hash [sha256.Size]byte) bool {
	return hash[0]%2 == 1
}
func SetFirstBit(hash [sha256.Size]byte, to bool) [sha256.Size]byte {
	hash[0] -= hash[0]%2  // Always even.
	if to {
		hash[0] += 1
	}
	return hash
}

//Copies it.. because no go stuff for that ><
func to_byte256(x []byte) [sha256.Size]byte {
	//assert len(x) <= sha256.size
	var ret [sha256.Size]byte
	i := 0
	for i < len(x) {
		ret[i] = x[i]
		i += 1
	}
	for i < sha256.Size {
		ret[i] = 0
		i += 1
	}
	return ret
}

// Too 'plain lengths' bytes, first bit zero.
func tbfb(x [sha256.Size]byte) []byte {
	ret := x[:]
	ret[0] -= ret[0]%2
	return ret
}

func H(a []byte) [sha256.Size]byte {
	return sha256.Sum256(a)
}

// NOTE: if you intend to maybe change, assume H_2(a,b) != H(b,a) is REQUIRED a-priori.
func H_2(a [sha256.Size]byte, b [sha256.Size]byte) [sha256.Size]byte {
	return SetFirstBit(sha256.Sum256(append(tbfb(a), tbfb(b)...)), false)
}

// A node of Merkle tree, note that the below omits a lot.
type MerkleNode struct {
	Hash   [sha256.Size]byte
	Left   *MerkleNode
	Right  *MerkleNode
	Up     *MerkleNode  // Otherwise it is messy to create the path.
}

func (self *MerkleNode) interest() bool {  // Even/odd is whether interest.
	return FirstBit(self.Hash)
}

func new_MerkleNode(left *MerkleNode,right *MerkleNode) *MerkleNode {
	hash := SetFirstBit(H_2(left.Hash, right.Hash), left.interest() || right.interest())
	node := &MerkleNode{Hash:hash, Left:left, Right:right, Up:nil}
	left.Up  = node
	right.Up = node
	return node
}

func selective_new_MerkleNode(left *MerkleNode,right *MerkleNode) *MerkleNode {
	node := new_MerkleNode(left, right)
	node.Left.effect_interest()
	node.Right.effect_interest()
	return node
}

func (node MerkleNode) effect_interest() {
	if !node.interest() { //No interest in branches.
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
	List []MerkleTreePortion
}

func NewMerkleTreeGen() *MerkleTreeGen {
	return &MerkleTreeGen{List:[]MerkleTreePortion{}}
}

// Adds chunk where you calculated the hash, returning the leaf the current is on.
func (gen *MerkleTreeGen) AddChunkH(h [sha256.Size]byte, interest bool) *MerkleNode {
	if len(gen.List) == 0 || gen.List[0].Depth != 1 {
		add_node := &MerkleNode{Hash:h, Left:nil, Right:nil, Up:nil}

		list := []MerkleTreePortion{}
		list = append(list, MerkleTreePortion{Node:add_node, Depth:1})
		gen.List = append(list, gen.List...)
		return add_node
	} else {
		// assert gen.List[0].Depth == 1
		new_leaf := &MerkleNode{Hash:h, Left:nil, Right:nil}

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
func (gen *MerkleTreeGen) AddChunk(chunk []byte, interest bool) *MerkleNode {
	return gen.AddChunkH(SetFirstBit(H(chunk), interest), interest)
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
	case node.Left != nil && node.Right != nil &&
		   H_2(node.Left.Hash, node.Right.Hash) != SetFirstBit(node.Hash, false):
		return node, false
	case recurse == 0 || node.Up == nil:
		return node, true
	default:
		return node.Up.IsValid(recurse - 1)
	}
}

func (node *MerkleNode) CorrespondsToChunk(chunk []byte) bool {
	return node.CorrespondsToHash(H(chunk))
}

func (node *MerkleNode) CorrespondsToHash(H [sha256.Size]byte) bool {
	return SetFirstBit(H, false) == SetFirstBit(node.Hash, false)
}

func (node* MerkleNode) Verify(Hroot [sha256.Size]byte, Hchunk [sha256.Size]byte) bool {
	root, internal := node.IsValid(-1)
	return internal && root.CorrespondsToHash(Hroot) && node.CorrespondsToHash(Hchunk)
}

// Calculated paths essentially make a compilation of the data needed to do the
// check. 
func (node *MerkleNode) Path() [][sha256.Size]byte {
	switch {
	case node.Right != nil || node.Left != nil:  return [][sha256.Size]byte{}
	case node.Up == nil:                       	 return [][sha256.Size]byte{}
	default:                                     return node.Up.path(node)
	}
}

func (node *MerkleNode) path(from *MerkleNode) [][sha256.Size]byte {
	path := [][sha256.Size]byte{}
	if node.Up != nil {
		path = node.Up.path(node)
	}

	switch {
	case node.Right == from:  return append(path, SetFirstBit(node.Left.Hash, true))
	case node.Left == from:   return append(path, SetFirstBit(node.Right.Hash, false))
	default:                  return [][sha256.Size]byte{} // Invalid Merkle tree.
	}
}

//Calculate expected root, given the path.
func ExpectedRoot(H_leaf [sha256.Size]byte, path [][sha256.Size]byte) [sha256.Size]byte {
	x := H_leaf
	for i := range path {
		h := path[len(path) - i - 1]
		switch {
		case FirstBit(h):  x = H_2(h, x)
		default:           x = H_2(x, h)
		}
	}
	return x
}

//Checks a root.
func Verify(root [sha256.Size]byte, leaf []byte, path [][sha256.Size]byte) bool {
	return SetFirstBit(ExpectedRoot(H(leaf), path), false) == SetFirstBit(root, false)
}
