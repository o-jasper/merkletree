package trie_easy

import "fmt"

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
	
	//TODO merkle-tree-like stuff.
}

// ----

type TrieNode struct {
	Actual TrieNodeInterface
}

func NewTrieNode(actual TrieNodeInterface) TrieNode { return TrieNode{Actual:actual} }

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

func (n* TrieNode) SetI(str []byte, j int64, to interface{}) {
	node, i := n.Downward(str, j)
	if node.Actual == nil { node.Actual = NewTrieNode16(nil) }
  node.Actual = node.Actual.SetRaw(str, i, to)
}

func (n* TrieNode) Set(str []byte, to interface{}) { n.SetI(str, 0, to) }

// ---- Plain 16-way split with data.

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
	//Make more trie nodes(TODO use TrieStretch)
/*	cur := n
	for i < 2*int64(len(str)) - 1 {
		m := NewTrieNode16(nil)
		cur.Sub[nibble(str,i)].Actual = m
		cur = m
		i += 1
	}
	cur.Sub[nibble(str,i)].Actual = &TrieNodeData{Data:to} //NewTrieNode16(to) */
	m := n
	if i%2 == 0 { // This one is even, so the next one isnt!
		m = NewTrieNode16(nil)
		n.Sub[str[i/2]%16] = NewTrieNode(m)
		i += 1
		if i == 2*int64(len(str)) { // Right, that fixed it.
			m.Data = to
			return n
		}
	}
	end := NewTrieNode(&TrieNodeData{Data:to})
	m.Sub[str[i/2]/16] = NewTrieNode(&TrieStretch{Stretch : str[i/2 + 1:], End : end})
	return n
}

// --- Just the data.
type TrieNodeData struct {
	Data interface{}
}

func (n *TrieNodeData) Downward(str []byte, i int64) (*TrieNode, int64) {
	return nil, i
}

func (n *TrieNodeData) Get(str []byte, i int64) interface{} {
	if i != 2*int64(len(str)) { panic("Endpoint only") }
	return n.Data
}

func (n* TrieNodeData) SetRaw(str []byte, i int64, to interface{}) TrieNodeInterface {
	if i == 2*int64(len(str)) {
		n.Data = to
		return n
	}
	return NewTrieNode16(n.Data).SetRaw(str, i, to)
}

// --- Stretch with just one branch.

type TrieStretch struct {
	Stretch  []byte
	End      TrieNode
}

func (n *TrieStretch) Downward(str []byte, i int64) (*TrieNode, int64) {
	if i < 2*int64(len(str) - len(n.Stretch)) { // Range inside the stretch.
		return nil, i
	}
	for j := int64(0) ; j < int64(len(n.Stretch)) ; j++ {
		if str[i/2 + j] != n.Stretch[j] {  // Breaks out of the stretch.
			return nil, i  // Back to begining.
		}
	}
	// Continue at end.
	return n.End.Downward(str, i + 2*int64(len(n.Stretch)))
}

func (n *TrieStretch) Get(str []byte, i int64) interface{} {
	if i < 2*int64(len(str) - len(n.Stretch)) { // Range inside the stretch.(nothing in there)
		return nil
	}
	return n.End.Get(str, i + 2*int64(len(n.Stretch)))
}

func (n* TrieStretch) SetRaw(str []byte, i int64, to interface{}) TrieNodeInterface {
	if i%2 != 0 { panic("TrieStretches should start at whole bytes.") }
	// The hard part.
	for j := int64(0) ; i/2 + j < int64(len(str)) ; j++ {
		if str[i/2 + j] != n.Stretch[j] {  // Breaks out of the stretch.
			a1, a2 := str[i/2 + j]%16, str[i/2 + j]/16  //Added nibbles.
			g1, g2 := n.Stretch[j]%16, n.Stretch[j]/16  //Already existing nibble route.

			fmt.Println("F", n.Stretch, i, j)

			// Two nodes(stretches start even)
			first, second := NewTrieNode16(nil), NewTrieNode16(nil)
			// Connect them.
			first.Sub[g1] = NewTrieNode(second)
			// Connect to what is after.
			if j < 2*int64(len(str)) { // It is a bit of stretch.
				fmt.Println("-", n.Stretch[j+1:], g1+16*g2)
				after_stretch := &TrieStretch{ Stretch : n.Stretch[j+1:], End : n.End }
				second.Sub[g2] = NewTrieNode(after_stretch)
			} else{  // It is the current end.
				second.Sub[g2] = n.End
			}

			if a1 != g1 { // Breaks out of first one.
				first.Sub[a1].SetI(str, i + 2*j, to)
			} else if a2 != g2 { // Breaks out of first one.
				if a1 != g1 { panic("BUG") }  // (could be both, first goes)
				second.Sub[a2].SetI(str, i + 2*j + 1, to)
			} else { panic("BUG") }
			
			if j == 0 {
				return first
			} else {  // Prepend what was before.
				fmt.Println("-", n.Stretch[:j])
				return &TrieStretch{Stretch : n.Stretch[:j], End : NewTrieNode(first)}
			}
			panic("Unreachable")
		}
	}
	panic("BUG Didnt go downward properly.(if it split it should also have cut out now.)")
	return nil
}
