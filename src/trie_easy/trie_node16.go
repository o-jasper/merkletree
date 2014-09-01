package trie_easy

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
	if i < 2*int64(len(str)) { return nil }
	if i > 2*int64(len(str)) { panic("i>len") }
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
	//Make more trie nodes(TODO use TrieStretch)
	
	final := &TrieNodeData{Data:to} //NewTrieNode16(to)
	n.Sub[nibble(str,i)].Actual = Creator16{}.Extend(str, i+1, final)
	return n
}

func (n* TrieNode16) MapAll(data interface{}, pre []byte, odd bool, fun MapFun) bool {
	if n.Data != nil && fun(data, pre, n.Data) { return true }
	for i, sub := range n.Sub {
		if sub.Actual == nil { continue }
		var npre []byte
		if odd {
			npre = append(pre[:len(pre)-1], pre[len(pre)-1] + byte(i)*16)
		} else {
			npre = append(pre, byte(i))
		}
		if sub.Actual.(TrieNodeInterface).MapAll(data, npre, !odd, fun){ return true }
	}
	return false
}

// -- Creates it that way.
type Creator16 struct {}

func (_ Creator16) Extend(str []byte, i int64, final TrieNodeInterface) interface{} {
	first := NewTrieNode16(nil)
	m := first
	for i < 2*int64(len(str)) - 1 {
		n := NewTrieNode16(nil)
		m.Sub[nibble(str,i)].Actual = n
		m = n
		i += 1
	}
	m.Sub[nibble(str,i)].Actual = final
	return first
}

// --- Just the data.
type TrieNodeData struct {
	Data interface{}
}

func (n *TrieNodeData) Downward(str []byte, i int64) (*TrieNode, int64) {
	return nil, i
}

func (n *TrieNodeData) Get(str []byte, i int64) interface{} {
	if i < 2*int64(len(str)) { return nil }
	if i > 2*int64(len(str)) { panic("i>2*len") }
	return n.Data
}

func (n* TrieNodeData) SetRaw(str []byte, i int64, to interface{}) TrieNodeInterface {
	if i == 2*int64(len(str)) {
		n.Data = to
		return n
	}
	return NewTrieNode16(n.Data).SetRaw(str, i, to)
}

func (n* TrieNodeData) MapAll(data interface{}, pre []byte, odd bool, fun MapFun) bool {
	return fun(data, pre, n.Data)	
}
