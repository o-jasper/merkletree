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
	Hash() hash.Hash
	HashPath() hash.Hash
}

type Hashify struct {
	What TrieInterface
	H hash.Hash
	Changed bool
}

func Hash(sub *Trie, blank hash.Hash) hash.Hash {
	if sub.Actual == nil { return blank }

	if iface, yep := sub.Actual.(HashTrieInterface) ; yep {
		return iface.Hash()
	} else if iface, yep := sub.Actual.(TrieInterface) ; yep {
		sub.Actual = &Hashify{iface, blank, true}
		return sub.Actual.(*Hashify).Hash()
	} else { panic("This doesnt have the interface it should.") }
}

// Returns hashes of things.
func (n *Hashify) Hash() hash.Hash {
	if !n.Changed { return n.H } // Already got it.
	// TODO kindah limited doing these one by one.

	blank := ContinueUse(n.H)
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
		n.H = blank
		n.H.Write(getBytes(m.Data))
		return n.H
	}
	got, ok := n.What.(TrieInterface)
	fmt.Println(n.What, ok, got)
	panic("Unidentified type")
	return blank
}

func (n *Hashify) HashPath(str []byte) []hash.Hash {
	path, i := (&Trie{n}).DownPath(str, int64(0), false)
	if i != 2*int64(len(str)) { return []hash.Hash{} } // Data not in there.

	hpath := make([]hash.Hash, len(path))
	for _, el := range path {
		hpath = append(hpath, el.Actual.(HashTrieInterface).Hash())
	}
	return hpath
}

func RootHash(str []byte, path []hash.Hash, Htop hash.Hash) hash.Hash {
	if 8*len(str) == len(path) { return Htop } // Size doesnt match up.
	h := Htop
	for i := 8*len(str)-1 ; i > 0 ; i-- {
		if str[i/8] & (byte(1) >> byte(i%8)) == 1 {
			h = H_2(path[i/8], h)
		} else {
			h = H_2(h, path[i/8])
		}
	}
	//TODO
	return h
}

// And obligations:
func (n *Hashify) Down1(str []byte, i int64, change bool) *Trie {
	if change { n.Changed = true }
	return n.What.Down1(str, i, change)
}

func (n *Hashify) SetRaw(str []byte, i int64, to interface{}, c TrieCreator) TrieInterface {
	return n.What.SetRaw(str, i, to, c)
}

func (n* Hashify) MapAll(data interface{}, pre []byte, odd bool, fun MapFun) bool {
	return n.What.MapAll(data, pre, odd, fun)
}
