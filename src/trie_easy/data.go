package trie_easy

// --- Just the data.
type DataNode struct { Data interface{} }

func (n DataNode) Down1(_ []byte, _ int64, _ bool) *Trie {
	return nil
}

func (n DataNode) Get(str []byte, i int64) interface{} {
	if i < 2*int64(len(str)) { return nil }
	if i > 2*int64(len(str)) { panic("i>2*len") }
	return n.Data
}

func (n DataNode) SetRaw(str []byte, i int64, to interface{}, c TrieCreator) TrieInterface {
	if i == 2*int64(len(str)) {
		n.Data = to
		return n
	}
	return NewNode16(n.Data).SetRaw(str, i, to, c)
}

func (n DataNode) MapAll(data interface{}, pre []byte, odd bool, fun MapFun) bool {
	return fun(data, pre, n.Data)	
}
