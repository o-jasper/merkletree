package trie_easy

func nibble(arr []byte, i int64) byte {
	if i % 2 == 0 {
		return arr[i/2] % 16
	}
	return arr[i/2] / 16
}

type TrieNodeInterface interface {
	Downward([]byte, int64) (*TrieNode, int64)  // Gets new trienodes insofar possible.
	Get([]byte, int64) interface{}
	SetRaw([]byte, int64, interface{}) TrieNodeInterface
}

// ----

type TrieNode struct {
	Actual TrieNodeInterface
}

func NewTrieNode() TrieNode {
	return TrieNode{Actual : nil}
}

func (n *TrieNode) Downward(str []byte, i int64) (*TrieNode, int64) {
	if n.Actual == nil { return n, i }
	if i == 2*int64(len(str)) {	return n, i }
	m, j := n.Actual.Downward(str, i)
	if i == j || m ==nil {	return n, i }  // Keep the last indirection.
	return m, j
}

func (n* TrieNode) Get(str []byte, i int64) interface{} {
	at, j := n.Downward(str, i)
	if at.Actual == nil {
		return nil
	}
	return at.Actual.Get(str, j)
}

func (n* TrieNode) Set(str []byte, to interface{}) {
	node, i := n.Downward(str, 0)
	if node.Actual == nil { node.Actual = NewTrieNode16(nil) }
  node.Actual = node.Actual.SetRaw(str, i, to)
}

// ----

type TrieNode16 struct {
	Sub  [16]TrieNode
	Data interface{}
}

func NewTrieNode16(data interface{}) *TrieNode16 {
	t := TrieNode16{Data:data}
	for i := 0 ; i < 16 ; i++ {
		t.Sub[i].Actual = nil
	}
	return &t
}

func (n *TrieNode16) Downward(str []byte, i int64) (*TrieNode, int64) {
	if i == 2*int64(len(str)) { panic("TrieNode16 `Downward` only for if actually need to.") }
	cur := &n.Sub[nibble(str, i)]
	i += 1
	for i < int64(2*len(str)) {
		if got, ok := cur.Actual.(*TrieNode16) ; ok {
			cur = &got.Sub[nibble(str, i)]
		} else {
			return cur.Downward(str, i)
		}
		i += 1
	}
	return cur, i
}

func (n *TrieNode16) Get(str []byte, i int64) interface{} {
	if i != 2*int64(len(str)) { panic("Endpoint only") }
	return n.Data
}

func (n* TrieNode16) SetRaw(str []byte, i int64, to interface{}) TrieNodeInterface {
	if i == 2*int64(len(str)) {
		n.Data = to
		return n
	}
	if n.Sub[nibble(str,i)].Actual != nil {
		panic("Not far downward enough!?!")
	}
	//Make more trie nodes(TODO special structure for sparse stuff.)
	cur := n
	for i < 2*int64(len(str)) - 1 {
		m := NewTrieNode16(nil)
		cur.Sub[nibble(str,i)].Actual = m
		cur = m
		i += 1
	}
	cur.Sub[nibble(str,i)].Actual = NewTrieNode16(to)
	return n
}
