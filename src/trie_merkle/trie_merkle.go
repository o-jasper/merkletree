package trie_merkle

import (
	"trie_easy"
	"hash"
//	"fmt"
)

type HashTrieInterface interface {
	trie_easy.TrieInterface
	Hash(func() hash.Hash) hash.Hash
}

type Hashify struct {
	What trie_easy.TrieInterface
	H hash.Hash
}

func Hash(of trie_easy.TrieInterface, blank func() hash.Hash) hash.Hash {
	if got, is_hash := of.(HashTrieInterface) ; is_hash {
		return got.Hash(blank)  //Already of correct interface.
	} else { // Need to give it an interface first.
		return (&Hashify{of, nil}).Hash(blank)
	}
}
//TODO need paths.

// Returns hashes of things.
func (n *Hashify) Hash(blank func() hash.Hash) hash.Hash {
	if n.H != nil { return n.H } // Already got it.
	n.H = blank()
	h_nothing := blank().Sum([]byte{})

	// TODO kindah limited doing these one by one.

	if m, ok := n.What.(*trie_easy.Node16) ; ok {
		for _, sub := range m.Sub { //TODO not acceptable.(15 per time instead of 4)
			if sub.Actual == nil {
				n.H.Write(h_nothing)
			} else if iface, yep := sub.Actual.(HashTrieInterface) ; yep {
				n.H.Write(iface.Hash(blank).Sum([]byte{}))
			} else if iface, yep := sub.Actual.(trie_easy.TrieInterface) ; yep {
				sub.Actual = Hashify{iface, nil}
			} else { panic("This doesnt have the interface it should.") }
		}
		n.H.Write(getBytes(m.Data))
		return n.H
	}
	panic("Unidentified type")
	return blank()
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
