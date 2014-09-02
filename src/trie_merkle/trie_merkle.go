package trie_merkle

import (
	"trie_easy"
	"hash"
)

type HashTrieInterface interface {
	trie_easy.TrieInterface
	Hash(func() hash.Hash) hash.Hash
}

type Hashify struct {
	What trie_easy.TrieNode
	H Hash	
}

func (n *Hashify) Hash(blank func() hash.Hash) hash.Hash {
	if n.H != nil { return n.H } // Already got it.
	n.H = blank()
	h_nothing := blank().Sum([]byte{})

	// TODO kindah limited doing these one by one.

	if m, ok := n.What.(trie_easyNode16) ; ok {
		for i, sub := range n.Sub {
			if sub.Actual == nil {
				n.H.Write(h_nothing)
			} else {
				got, ok := sub.Actual.(HashTrieInterface)
				if !ok {
					sub.Actual = Hashify{sub.actual, nil}
					got = sub.Actual
				}
				n.H.Write(got.Hash(blank).Sum([]byte{}))
			}
		}
		n.H.Write(getBytes(n.Data))
		return n.H
	}
}

func (n *Hashify) Down1(str []byte, i int64, change bool) *Trie {
	if change { n.H = nil }
	return n.What.Down1(str, i, change)
}

func (n *Hashify) SetRaw(str []byte, i int64, to interface{}, c TrieCreator) TrieInterface {
	return n.What.SetRaw(str, i, to, c)
}

func (n* ) MapAll(data interface{}, pre []byte, odd bool, fun MapFun) bool {
