package trie_easy

// The interface and stuff using the interface.

type MapFun func(interface{}, []byte, interface{}) bool

type TrieInterface interface {
	Downward([]byte, int64) (*Trie, int64)  // Gets new trienodes insofar possible.
	Get([]byte, int64) interface{}
	SetRaw([]byte, int64, interface{}) TrieInterface

	MapAll(interface{}, []byte, bool, MapFun) bool

	//TODO merkle-tree-like stuff.
}

// Allows implementation of different ways to extend it.
type TrieCreator interface {
	Extend([]byte, int64, TrieInterface) interface{}
}

// ----

type Trie struct {
	Actual interface{}  // Its an interface for ability to expand.
}

func NewTrie(actual TrieInterface) Trie { return Trie{Actual:actual} }

func (n *Trie) Downward(str []byte, i int64) (*Trie, int64) {
	if n.Actual == nil { return n, i }
	if i == 2*int64(len(str)) {	return n, i }
	m, j := n.Actual.(TrieInterface).Downward(str, i)
	if i == j || m ==nil {	return n, i }  // Keep the last indirection.
	return m, j
}

func (n* Trie) Get(str []byte, i int64) interface{} {
	at, j := n.Downward(str, i)
	if at.Actual == nil {	return nil }
	return at.Actual.(TrieInterface).Get(str, j)
}

func (n* Trie) SetI(str []byte, j int64, to interface{}) {
	node, i := n.Downward(str, j)
	if node.Actual == nil { node.Actual = NewNode16(nil) }
  node.Actual = node.Actual.(TrieInterface).SetRaw(str, i, to)
}

func (n* Trie) Set(str []byte, to interface{}) { n.SetI(str, 0, to) }

func (n* Trie) MapAll(data interface{}, fun MapFun) bool {
	if n.Actual == nil { return false }
	return n.Actual.(TrieInterface).MapAll(data, []byte{}, false, fun)
}
