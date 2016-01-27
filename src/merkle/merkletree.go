
package merkle

import . "hash_extra"

// A node of Merkle tree, note that the below omits a lot.
type MerkleNode struct {
	Hash   HashResult
	Left   *MerkleNode
	Right  *MerkleNode
	Up     *MerkleNode  // Otherwise it is messy to create the path.
	Interest bool
}

func selective_new_MerkleNode(h Hasher, left, right *MerkleNode) *MerkleNode {
	n := MerkleNode{Hash : h.H_U2(left.Hash, right.Hash), Left : left, Right : right, Up:nil}
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
	Hasher Hasher
	I uint64
	List []MerkleTreePortion
	// Whether to prepend the index in the chunks.(otherwise order is not recorded at all)
	IncludeIndex bool
}

func NewMerkleTreeGen(hasher Hasher, include_index bool) MerkleTreeGen {
	return MerkleTreeGen{hasher, uint64(0), []MerkleTreePortion{}, include_index}
}

// Adds chunk where you calculated the hash, returning the leaf the current is on.
func (gen *MerkleTreeGen) AddH(h HashResult, interest bool) *MerkleNode {
	//if len(h) != gen.Hasher.Size() { panic("Does not look like the correct hash.") }
	gen.I += 1
	if len(gen.List) == 0 || gen.List[0].Depth != 1 {
		add_node := &MerkleNode{Hash:h, Left:nil, Right:nil, Up:nil, Interest:interest}

		list := []MerkleTreePortion{}
		list = append(list, MerkleTreePortion{Node:add_node, Depth:1})
		gen.List = append(list, gen.List...)
		return add_node
	} else {
		if gen.List[0].Depth != 1 { panic("Depth of first element should be 1") }
		new_leaf := &MerkleNode{Hash:h, Left:nil, Right:nil, Interest:interest}
		//Combine the two.
		new_node := selective_new_MerkleNode(gen.Hasher, new_leaf, gen.List[0].Node)
		gen.List[0] = MerkleTreePortion{Node:new_node, Depth:2}

		// Combine more, while equal depth.
		for len(gen.List) >= 2 && gen.List[1].Depth == gen.List[0].Depth {
			new_node := selective_new_MerkleNode(gen.Hasher, gen.List[0].Node, gen.List[1].Node)
			gen.List = gen.List[1:] // Cut off the first one.
			gen.List[0] = MerkleTreePortion{Node:new_node, Depth:gen.List[0].Depth + 1}
		}
		return new_leaf  //Return the leaf.
	}
}

func (gen *MerkleTreeGen) Add(leaf []byte, interest bool) *MerkleNode {
	if gen.IncludeIndex {
		return gen.AddH(gen.Hasher.HwI(gen.I, leaf), interest)
	} else {
		return gen.AddH(gen.Hasher.H(leaf), interest)
	}
}

// Coerce the last parts together, returning the root.
// NOTE: you can 'finish' and then continue to make what you put in already
// becomes a bit of a 'lob' that takes longer Merkle paths.
func (gen *MerkleTreeGen) Finish() *MerkleNode {
	// assert len(gen.List) > 0
	for len(gen.List) >= 2  {
		new_node := selective_new_MerkleNode(gen.Hasher, gen.List[0].Node, gen.List[1].Node)
		gen.List = gen.List[1:]
		gen.List[0] = MerkleTreePortion{Node:new_node, Depth:gen.List[0].Depth}
	}
	return gen.List[0].Node
}

// Only checks internal consistency upward;
// Checks a merkle path, _except_ the root and leaf.
func (node *MerkleNode) IsValid(hasher Hasher, recurse int32) (*MerkleNode, bool) {
	cur := node
	for recurse != 0 {
		if cur.Left != nil && cur.Right != nil &&	!HashEqual(hasher.H_U2(cur.Left.Hash, cur.Right.Hash), cur.Hash) {
			return cur, false
		}
		recurse -= 1
		if cur.Up == nil { return cur, true }
		cur = cur.Up
	}
	return cur, true
}

func (node *MerkleNode) CorrespondsH(h HashResult) bool {
	return ByteSliceEqual(h, node.Hash)
}
func (node *MerkleNode) Corresponds(hasher Hasher, leaf []byte) bool {
	return node.CorrespondsH(hasher.H(leaf))
}
func (node *MerkleNode) CorrespondsWithIndex(hasher Hasher, i uint64, leaf []byte) bool {
	return node.CorrespondsH(hasher.HwI(i, leaf))
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

func (node* MerkleNode) VerifyH(hasher Hasher, Hroot, Hleaf HashResult) int8 {
	root, internal := node.IsValid(hasher, -1)
	switch {
	case !internal:                  return WrongDataPath
	case !node.CorrespondsH(Hleaf):  return WrongDataLeaf
	case !root.CorrespondsH(Hroot):  return WrongDataRoot
	default:                         return Correct
	}
}

func (node* MerkleNode) Verify(hasher Hasher, Hroot HashResult, leaf []byte) int8 {
	return node.VerifyH(hasher, Hroot, hasher.H(leaf))
}

func (node* MerkleNode) VerifyWithIndex(hasher Hasher, Hroot HashResult, i uint64, leaf []byte) int8 {
	return node.VerifyH(hasher, Hroot, hasher.HwI(i, leaf))
}

// Calculated paths essentially make a compilation of the data needed to do the
// check. 
func (node *MerkleNode) Path() []HashResult {
	prev, cur, list := node, node, []HashResult{}
	for cur != nil {
		if cur.Left  == prev {
			list = append(list, cur.Right.Hash)
			if cur.Right == prev { panic("Both?") }
		}	else if cur.Right == prev { list = append(list, cur.Left.Hash) }
		prev = cur
		cur = cur.Up
	}
	return list
}
