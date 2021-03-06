package trie_easy

type Node16 struct {
	Sub  [16]Trie
	Data interface{}
}

func NewNode16(data interface{}) *Node16 {
	t := Node16{Data:data}
	for i := 0 ; i < 16 ; i++ { t.Sub[i].Actual = nil }
	return &t
}

func (n *Node16) Down1(str []byte, i int64, _ bool) *Trie {
	return &n.Sub[Nibble(str, i)]
}

func (n *Node16) Get(str []byte, i int64) interface{} {
	if i < 2*int64(len(str)) { return nil }
	if i > 2*int64(len(str)) { panic("i>len") }
	return n.Data
}

func (n* Node16) SetRaw(str []byte, i int64, to interface{}, c TrieCreator) TrieInterface {
	if i == 2*int64(len(str)) {
		n.Data = to
		return n
	}
	if n.Sub[Nibble(str,i)].Actual != nil {
		panic("Not far downward enough!?!")
	}
	n.Sub[Nibble(str,i)].Actual = c.Extend(str, i+1, DataNode{Data:to})
	return n
}

func (n* Node16) MapAll(data interface{}, pre []byte, odd bool, fun MapFun) bool {
	if n.Data != nil && fun(data, pre, n.Data) { return true }
	for i, sub := range n.Sub {
		if sub.Actual == nil { continue }
		var npre []byte
		if odd {
			npre = append(pre[:len(pre)-1], pre[len(pre)-1] + byte(i)*16)
		} else {
			npre = append(pre, byte(i))
		}
		if sub.Actual.(TrieInterface).MapAll(data, npre, !odd, fun){ return true }
	}
	return false
}

// -- Creates it that way.
type Creator16 struct {}

func (_ Creator16) Extend(str []byte, i int64, final TrieInterface) interface{} {
	first := NewNode16(nil)
	m := first
	for i < 2*int64(len(str)) - 1 {
		n := NewNode16(nil)
		m.Sub[Nibble(str,i)].Actual = n
		m = n
		i += 1
	}
	m.Sub[Nibble(str,i)].Actual = final
	return first
}
