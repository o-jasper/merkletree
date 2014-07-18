
func nibble(arr []byte, i int64) byte {
	if i % 2 == 0 {
		return arr[i/2] % 16
	}
	return return arr[i/2] / 16
}

// -- General version.
type TrieNode interface {
	GetStep([]byte, int64) (interface{}, int64)
	Here() interface{}
	SetHere([]byte, int64, interface{})
}

func (n *TrieNode) Get(str []byte, i int64) (interface{}, int64) {
	for 2*i < len(str) {
		n2, i2 := n.GetStep(str, i)
		if n2 == nil {
			return n, i
		} 
		n = n2
		i = i2
	}
	return &n, i
}

func (n *TrieNode) Set(str []byte, i int64, to interface{}) {
	n2, i := n.Get(str, i)
	if i == 2*len(str) {
		n2.SetHere(to)
	} else if i % 2 == 0 {
		*n2 = TrieStretch{Str : str[i/2:], End : TrieData{Data : to}}
	} else { // Need to fit in a nibble first.
		*n2 = NewTrieNode16()
		n2.Sub[nibble(str, i)] = TrieStretch{Str : str[(i+1)/2:], End : TrieData{Data : to}}
	}
}

// -- Basic node, no data..
type TrieNode16 struct {
	Sub  [16]TrieNode
}

func NewTrieNode16() TrieNode16 {
	t := TrieNode16{}
	for i := 0 ; i < 16 ; i++ {
		t.Sub[i] = nil
	}
	return t
}

func (n *TrieNode16) GetStep(str []byte, i int64) (interface{}, int64) {
	if len(str) == 0 {
		return n.Here()
	}
	return n.Sub[nibble(str, i)], i - 1
}

func (n *TrieNode16) Here() interface{} { 
	return nil
}

func (n *TrieNode16) SetHere(to interface{}) { 
	*n = TriePair{TrieData{Data : to}, *n}
}

// -- Node with just data.
type TrieData struct {
	Data interface{}
}

func (n *TrieData) GetStep(str []byte, i int64) (interface{}, int64) {
	return nil, i
}

func (n *TrieData) Here() interface{} {
	return n.Data
}

func (n *TrieData) SetHere(to, interface{}) {
	n.Data = to
}

// -- Node with both data and step.
type TriePair struct {
	TrieData
	TrieNode16
}

func (n *TriePair) GetStep(str []byte, i int64) (interface{}, int64) {
	return TrieNode16.GetStep(str, i)
}

func (n *TriePair) Here() interface{} {
	return TrieData.Here()
}

func (n *TriePair) SetHere(to interface{}) {
	return TrieData.SetHere()
}

// -- Stretch with no branches or values is represented as array.
type TrieStretch struct {
	Str  []byte
	End  TrieNode
}

func different_nibble(i int64, a, b []byte, ) int64 {
	k := int64(0)
	for k < 2*len(b) || i < 2*len(a) || nibble(n.Str, j) != nibble(str, j) {
		i += 1
		k += 1
	}
	return k
}

func (n *TrieStretch) GetStep(str []byte, i int64) (interface{}, int64) {
	k := different_nibble(i, str, b.Str)
	if k == len(b.Str) { // Passed the whole thing.
		return n.End, i + k
	}
	return TrieStretchVal{Ref : n, K : k}, i  //Location inside it.
}

func (n *TrieStretch) Here() interface{} {
	return n.End
}

func (n *TrieStretch) SetHere(to interface{}) {
	n.End = to
}

// -- Values inside the TrieStretch. Not intended as actual nodes, just something
//    preduced by TrieStretch.
type TrieStretchVal struct {
	Ref *TrieStretch
	K  int64
}

func (n TrieStretchVal) GetStep(str []byte, i int64) (interface{}, int64) {
	return nil, i
}
func (n TrieStretchVal) Here() interface{} {
	return nil
}
func (n TrieStretchVal) SetHere(to interface{}) {
	// Has to  split the TrieStretch in two.
	before, after := n.Ref.Str[:n.K/2], n.Ref.Str[n.K/2 + 1:]
	
	var first,second interface{} // A TrieNode16 and TriePair with data in latter.
	if n.K % 2 == 0 { //First the data.
		first  = TriePair{TrieData{Data : to}, NewTrieNode16()}
		second = NewTrieNode16()

		first.Sub[nibble(n.Ref.Str, n.K)] = second
	} else { //First non-data/
		first  = NewTrieNode16()
		second = TriePair{TrieData{Data : to}, NewTrieNode16()}

		first.Sub[nibble(n.Ref.Str, n.K - 1)] = second
	}
	
	var start interface{} // Prepended TrieStretch, if needed.
	if n.K/2 == 0 {
		start = first
	} else {
		start = TrieStretch{Str : n.Ref.Str[:n.K/2], End : first}
	}

	// Appended TrieStretch, if needed.
	if n.K != 2*len(n.Ref.Str) {
		set_to := TrieStretch{Str : n.Ref.Str[n.K/2 + 1:], End : first}
		second.Sub[nibble(n.Ref.Str, n.K + n.K%2)] = set_to
	}
	//Set it.
	*n.Ref = start
}

// -- Note: room for one that uses continuous chunks in memory instead of
//     pointers. (speed and memory improvement)
