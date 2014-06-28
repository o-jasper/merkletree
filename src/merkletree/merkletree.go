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
	//assert len(x)==sha256.size
	var ret [sha256.Size]byte
	copy(ret[:], x)
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

// Adds chunk, returning the leaf the current is on
func (gen *MerkleTreeGen) AddChunk(chunk []byte, interest bool) *MerkleNode {
	h := SetFirstBit(H(chunk), interest)
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

// Self-check. Not that since you can remove the rest of the nodes, this
// essentially already does proving by paths-and-chunks.
func (node *MerkleNode) IsValid(recurse int32) bool {
	if node.Left != nil && node.Right != nil {
		if H_2(node.Left.Hash, node.Right.Hash) != SetFirstBit(node.Hash, false) {
			return false
		}
	}
	if recurse == 0 || node.Up == nil {
		return true
	}
	return node.Up.IsValid(recurse - 1)
}

func (node *MerkleNode) CorrespondsToChunk(chunk []byte) bool {
	return SetFirstBit(H(chunk), false) == SetFirstBit(node.Hash, false) && 
		     node.Right == nil && node.Left == nil
}

// Calculated paths essentially make a compilation of the data needed to do the
// check. 
func (node *MerkleNode) Path() [][sha256.Size]byte {
	if node.Right != nil || node.Left != nil {
		return [][sha256.Size]byte{}
	} else if node.Up == nil {
		return [][sha256.Size]byte{}
	} else {
		return node.Up.path(node)
	}
}

func (node *MerkleNode) path(from *MerkleNode) [][sha256.Size]byte {
	path := [][sha256.Size]byte{}
	if node.Up != nil {
		path = node.Up.path(node)
	}

	if node.Right == from { //Came from right.
		return append(path, SetFirstBit(node.Left.Hash, true))
	}	else if node.Left == from { //Came from left.
		return append(path, SetFirstBit(node.Right.Hash, false))
	}	else { // Information was not stored, or invalid Merkle tree.
		return [][sha256.Size]byte{}
	}
}

//Calculate expected root, given the path.
func ExpectedRoot(H_leaf [sha256.Size]byte, path [][sha256.Size]byte) [sha256.Size]byte {
	x := H_leaf
	for i := range path {
		h := path[len(path) - i - 1]
		if FirstBit(h) {
			x = H_2(h, x)
		} else {
			x = H_2(x, h)
		}
	}
	return x
}

//Checks a root.
func CorrectRoot(root [sha256.Size]byte, leaf []byte, path [][sha256.Size]byte) bool {
	return SetFirstBit(ExpectedRoot(H(leaf), path), false) == SetFirstBit(root, false)
}
