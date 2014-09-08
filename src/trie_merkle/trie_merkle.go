package trie_merkle

import (
	"trie_easy"
	"hash"
	"fmt"
)

func H_2(a, b hash.Hash, blank func() hash.Hash) hash.Hash {
	h := blank()
	h.Write(a.Sum([]byte{}))
	h.Write(b.Sum([]byte{}))
	return h
}

//TODO drat.. Needs to be 'indexed'..?..

type HashTrieInterface interface {
	trie_easy.TrieInterface
	Hash(func() hash.Hash) hash.Hash
	HashPath(func() hash.Hash) hash.Hash
}

type Hashify struct {
	What trie_easy.TrieInterface
	H hash.Hash
}

func Hash(sub *trie_easy.Trie, blank func() hash.Hash) hash.Hash {
	if sub.Actual == nil { return blank() }

	if iface, yep := sub.Actual.(HashTrieInterface) ; yep {
		return iface.Hash(blank)
	} else if iface, yep := sub.Actual.(trie_easy.TrieInterface) ; yep {
		sub.Actual = &Hashify{iface, nil}
		return sub.Actual.(*Hashify).Hash(blank)
	} else { panic("This doesnt have the interface it should.") }
}

// Returns hashes of things.
func (n *Hashify) Hash(blank func() hash.Hash) hash.Hash {
	if n.H != nil { return n.H } // Already got it.
	// TODO kindah limited doing these one by one.

	if n.What == nil { panic("Hashify may not have nil item") }
	if m, ok := n.What.(*trie_easy.Node16) ; ok {
		h := make([]hash.Hash, 8, 8)  // Note: can be done with less memory.
		for i := 0 ; i < 16 ; i += 2 {
			h[i/2] = H_2(Hash(&m.Sub[i], blank), Hash(&m.Sub[i+1], blank), blank)
		}
		for i := 0 ; i < 8; i += 2   { h[i/2] = H_2(h[i], h[i+1], blank) }
		for i := 0 ; i < 4 ; i += 2  { h[i/2] = H_2(h[i], h[i+1], blank) }
		n.H	= H_2(h[0], h[1], blank)
		n.H.Write(getBytes(m.Data))
		return n.H
	}
	if m, ok := n.What.(trie_easy.DataNode) ; ok {
		n.H = blank()
		n.H.Write(getBytes(m.Data))
		return n.H
	}
	got, ok := n.What.(trie_easy.TrieInterface)
	fmt.Println(n.What, ok, got)
	panic("Unidentified type")
	return blank()
}

func (n *Hashify) HashPath(str []byte, blank func() hash.Hash) []hash.Hash {
	path, i := (&trie_easy.Trie{n}).DownPath(str, int64(0), false)
	if i != 2*int64(len(str)) { return []hash.Hash{} } // Data not in there.

	hpath := make([]hash.Hash, len(path))
	for _, el := range path {
		hpath = append(hpath, el.Actual.(HashTrieInterface).Hash(blank))
	}
	return hpath
}

func (n *Hashify) Down1(str []byte, i int64, change bool) *trie_easy.Trie {
	if change { n.H = nil }
	return n.What.Down1(str, i, change)
}

func (n *Hashify) SetRaw(str []byte, i int64, to interface{}, c trie_easy.TrieCreator) trie_easy.TrieInterface {
	return n.What.SetRaw(str, i, to, c)
}

func (n* Hashify) MapAll(data interface{}, pre []byte, odd bool, fun trie_easy.MapFun) bool {
	return n.What.MapAll(data, pre, odd, fun)
}
