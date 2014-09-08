package trie_merkle

import (
	"hash"
	"fmt"
)
import . "trie_easy"
import . "hash_extra"

//TODO drat.. Needs to be 'indexed'..?..

type HashTrieInterface interface {
	TrieInterface
	Hash(func() hash.Hash) hash.Hash
	HashPath(func() hash.Hash) hash.Hash
}

type Hashify struct {
	What TrieInterface
	H hash.Hash
}

func Hash(sub *Trie, blank func() hash.Hash) hash.Hash {
	if sub.Actual == nil { return blank() }

	if iface, yep := sub.Actual.(HashTrieInterface) ; yep {
		return iface.Hash(blank)
	} else if iface, yep := sub.Actual.(TrieInterface) ; yep {
		sub.Actual = &Hashify{iface, nil}
		return sub.Actual.(*Hashify).Hash(blank)
	} else { panic("This doesnt have the interface it should.") }
}

// Returns hashes of things.
func (n *Hashify) Hash(blank func() hash.Hash) hash.Hash {
	if n.H != nil { return n.H } // Already got it.
	// TODO kindah limited doing these one by one.

	if n.What == nil { panic("Hashify may not have nil item") }
	if m, ok := n.What.(*Node16) ; ok {
		h := make([]hash.Hash, 8, 8)  // Note: can be done with less memory.
		for i := 0 ; i < 16 ; i += 2 {
			h[i/2] = H_2(Hash(&m.Sub[i], blank), Hash(&m.Sub[i+1], blank))
		}
		for i := 0 ; i < 8; i += 2   { h[i/2] = H_2(h[i], h[i+1]) }
		for i := 0 ; i < 4 ; i += 2  { h[i/2] = H_2(h[i], h[i+1]) }
		n.H	= H_2(h[0], h[1])
		n.H.Write(getBytes(m.Data))
		return n.H
	}
	if m, ok := n.What.(DataNode) ; ok {
		n.H = blank()
		n.H.Write(getBytes(m.Data))
		return n.H
	}
	got, ok := n.What.(TrieInterface)
	fmt.Println(n.What, ok, got)
	panic("Unidentified type")
	return blank()
}

func (n *Hashify) HashPath(str []byte, blank func() hash.Hash) []hash.Hash {
	path, i := (&Trie{n}).DownPath(str, int64(0), false)
	if i != 2*int64(len(str)) { return []hash.Hash{} } // Data not in there.

	hpath := make([]hash.Hash, len(path))
	for _, el := range path {
		hpath = append(hpath, el.Actual.(HashTrieInterface).Hash(blank))
	}
	return hpath
}

func CheckHashPath(str []byte, path []hash.Hash, chunk interface{}) bool {
	//TODO
	return false
}

// And obligations:
func (n *Hashify) Down1(str []byte, i int64, change bool) *Trie {
	if change { n.H = nil }
	return n.What.Down1(str, i, change)
}

func (n *Hashify) SetRaw(str []byte, i int64, to interface{}, c TrieCreator) TrieInterface {
	return n.What.SetRaw(str, i, to, c)
}

func (n* Hashify) MapAll(data interface{}, pre []byte, odd bool, fun MapFun) bool {
	return n.What.MapAll(data, pre, odd, fun)
}
