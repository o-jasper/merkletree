package trie_easy

// The interface and stuff using the interface.

type MapFun func(interface{}, []byte, interface{}) bool

type TrieInterface interface {
	Down1([]byte, int64, bool) *Trie
	
	Get([]byte, int64) interface{}
	SetRaw([]byte, int64, interface{}, TrieCreator) TrieInterface

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

func (n *Trie) Downward(str []byte, i int64, changing bool) (*Trie, int64) {
	if n.Actual == nil { return n, i }
	m := n
	for i < 2*int64(len(str)) {
		if m.Actual == nil { return n, i - 1 }
		n = m
		iface, ok := m.Actual.(TrieInterface)
		if !ok { panic("Not trieinterfaable!!") }
		m = iface.Down1(str, i, changing)
		i += 1
	}
	return m, i
}

func (n* Trie) Get(str []byte, i int64) interface{} {
	at, j := n.Downward(str, i, false)
	if at.Actual == nil {	return nil }
	return at.Actual.(TrieInterface).Get(str, j)
}

func (n* Trie) SetI(str []byte, j int64, to interface{}, c TrieCreator) {
	node, i := n.Downward(str, j, true)
	if node.Actual == nil { node.Actual = NewNode16(nil) }
  node.Actual = node.Actual.(TrieInterface).SetRaw(str, i, to, c)
}

func (n* Trie) Set(str []byte, to interface{}, c TrieCreator) { 
	n.SetI(str, 0, to, c)
}

func (n* Trie) MapAll(data interface{}, fun MapFun) bool {
	if n.Actual == nil { return false }
	return n.Actual.(TrieInterface).MapAll(data, []byte{}, false, fun)
}
