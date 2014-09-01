package trie_easy

// The interface and stuff using the interface.

type MapFun func(interface{}, []byte, interface{}) bool

type TrieNodeInterface interface {
	Downward([]byte, int64) (*TrieNode, int64)  // Gets new trienodes insofar possible.
	Get([]byte, int64) interface{}
	SetRaw([]byte, int64, interface{}) TrieNodeInterface

	MapAll(interface{}, []byte, bool, MapFun) bool

	//TODO merkle-tree-like stuff.
}

// Allows implementation of different ways to extend it.
type TrieCreator interface {
	Extend([]byte, int64, TrieNodeInterface) interface{}
}

// ----

type TrieNode struct {
	Actual interface{}  // Its an interface for ability to expand.
}

func NewTrieNode(actual TrieNodeInterface) TrieNode { return TrieNode{Actual:actual} }

func (n *TrieNode) Downward(str []byte, i int64) (*TrieNode, int64) {
	if n.Actual == nil { return n, i }
	if i == 2*int64(len(str)) {	return n, i }
	m, j := n.Actual.(TrieNodeInterface).Downward(str, i)
	if i == j || m ==nil {	return n, i }  // Keep the last indirection.
	return m, j
}

func (n* TrieNode) Get(str []byte, i int64) interface{} {
	at, j := n.Downward(str, i)
	if at.Actual == nil {	return nil }
	return at.Actual.(TrieNodeInterface).Get(str, j)
}

func (n* TrieNode) SetI(str []byte, j int64, to interface{}) {
	node, i := n.Downward(str, j)
	if node.Actual == nil { node.Actual = NewTrieNode16(nil) }
  node.Actual = node.Actual.(TrieNodeInterface).SetRaw(str, i, to)
}

func (n* TrieNode) Set(str []byte, to interface{}) { n.SetI(str, 0, to) }

func (n* TrieNode) MapAll(data interface{}, fun MapFun) bool {
	if n.Actual == nil { return false }
	return n.Actual.(TrieNodeInterface).MapAll(data, []byte{}, false, fun)
}
