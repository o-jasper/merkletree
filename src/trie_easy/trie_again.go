package trie_easy

func nibble(arr []byte, i int64) byte {
	if i % 2 == 0 {
		return arr[i/2] % 16
	}
	return arr[i/2] / 16
}

type TrieNodeInterface interface {
	Downward([]byte, int64) (TrieNode, int64)  // Gets new trienodes insofar possible.
	Get([]byte, int64) interface{}
	SetRaw([]byte, int64, interface{}) TrieNodeInterface
}

// ----

type TrieNode struct {
	Actual TrieNodeInterface
}

func (n *TrieNode) Downward(str []byte, i int64) (TrieNodeInterface, int64) {
	m, j := n.Actual.Downward(str, i)
	if n.Actual == m {
		return n, j  // Keep the last indirection.
	}
	return m, j
}

func (n* TrieNode) Get(str []byte. i int64) interface{} {
	at, j := n.Downward(str, i)
	at.Actual.Get(str, i)
}

func (n* TrieNode) Set(str []byte. i int64, to interface{}) {
	n.Actual = n.Actual.SetRaw(str, i, to)
}

// ----

type TrieEnd struct {
	Sub  [16]TrieNode
}


// ----

type TrieNode16 struct {
	Sub  [16]TrieNode
}

func NewTrieNode16() TrieNode16 {
	t := TrieNode16{}
	for i := 0 ; i < 16 ; i++ {
		t.Sub[i].Actual= = nil
	}
	return t
}

func (n *TrieNode16) Downward(str []byte, i int64) (TrieNodeInterface, int64) {
	m := n
	for i < 2*len(str) {
		cur = n.Sub[nibble(str, i)].Actual
		i += 1
		if node, ok := cur.Actual.(*TrieNode16) ; !ok {
			return m.Downward(str, i)
		}
		m = n
	}
	return m
}

func (n *TrieNode16) Get(str []byte, i int64) interface{} {
	//if 2*i < int64(len(str)) { return nil }
	return nil
}

func (n* TrieNode16) SetRaw(str []byte. i int64, to interface{}) TrieNode {
	if 2*i == int64(len(str)) {
		return &TrieNodeBoth{*n, TrieNodeData{Data:to}}
	}
	//Make more trie nodes(TODO special structure for sparse stuff.)
	cur := n
	for 2*i < int64(len(str)) - 1 {
		m := NewTrieNode16()
		cur.Sub[nibble(str,i)].Actual = m
		cur = m
		i += 1
	}
	cur.Sub[nibble(str,i)].Actual = &TrieNodeBoth{NewTrieNode16, TrieNodeData{Data:to}}
	return n
}

// ----

type TrieNodeData struct {
	Data interface{}
}

func (n *TrieNodeData) Downward(str []byte, i int64) (TrieNodeInterface, int64) {
	return n, i
}

func (n *TrieNodeData) Get(str []byte, i int64) interface{} {
	if 2*i != int64(len(str)) { return nil }
	return n.Data
}

func (n* TrieNodeData) SetRaw(str []byte. i int64, to interface{}) TrieNode {
	if 2*i == int64(len(str)) {
		n.Data = to
		return n
	}
	m = TrieNodeBoth{NewTrieNode16(), *n}
	return m.TrieNode16.SetRaw(str, i ,to)
}

// ----

type TrieNodeBoth struct {
	TrieNode16
	TrieNodeData
}


func (n *TrieNodeData) Downward(str []byte, i int64) (TrieNodeInterface, int64) {
	return n.TrieNode16.Downward(str, i)
}

func (n *TrieNodeData) Get(str []byte, i int64) interface{} {
	return n.TrieNodeData.Get(str, i)
}

func (n* TrieNodeBoth) SetRaw(str []byte. i int64, to interface{}) TrieNode {
	if 2*i == int64(len(str)) {
		n.Data = to
		return n
	}
	n.TrieNode16.SetRaw(str, i ,to)
	return n
}
