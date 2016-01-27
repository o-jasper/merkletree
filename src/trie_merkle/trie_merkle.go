package trie_merkle

import (
	"hash"
//	"fmt"
)
import . "trie_easy"
import . "hash_extra"

//TODO drat.. Needs to be 'indexed'..?..

type HashTrieInterface interface {
	TrieInterface
	Hash() HashResult
	HashDown1(str []byte, i int64, list *[]HashResult) *Trie
}

func Hashify(input interface{}) HashTrieInterface {
	switch input.(type) {
	case *Node16:
		return &MT16{*input.(*Node16), false,
			[8]HashResult{}, [4]HashResult{}, [2]HashResult{}}
	case *DataNode:
		return &MTData{*input.(*DataNode)}
	}
	panic("")
	return nil
}

func Hash(sub *Trie, blank hash.Hash) hash.Hash {
	if sub.Actual == nil { return blank }

	if iface, yep := sub.Actual.(HashTrieInterface) ; yep {
		return iface.Hash()
	} else if iface, yep := sub.Actual.(TrieInterface) ; yep {
		sub.Actual = Hashify(iface, blank)
		return sub.Actual.(HashTrieInterface).Hash()
	} else { panic("This doesnt have the interface it should.") }
}

// Go down as far as possible, making a merkle path.
// NOTE: code pretty much the same as Downward
func HashPath(n *Trie, str []byte, blank hash.Hash) (*Trie, int64, []hash.Hash) {
	i := int64(0)
	if n.Actual == nil { return n, i, []hash.Hash{} }
	path, m := []hash.Hash{}, n
	for i < 2*int64(len(str)) {
		if m == nil || m.Actual == nil { return n, i - 1, path }
		n = m
		iface, ok := m.Actual.(HashTrieInterface)
		if !ok {
			Hash(m, blank)
			iface = m.Actual.(HashTrieInterface)
		}
		m = iface.HashDown1(str, i, &path)
		i += 1
	}
	return m, i, path
}
// Calculating, given data.

func RootHashH(str []byte, path []hash.Hash, Hleaf hash.Hash) hash.Hash {
	if 8*len(str) != len(path) { return Hleaf } // Size doesnt match up.
	h := Hleaf
	for i := 8*(len(str)-1) ; i > 0 ; i-- {
		if str[i/8] & (byte(1) >> byte(i%8)) == 1 {
			h = H_2(path[i/8], h)
		} else {
			h = H_2(h, path[i/8])
		}
	}
	return h
}
func VerifyH(str []byte, path []hash.Hash, Hleaf, root hash.Hash) bool {
	return HashEqual(RootHashH(str, path, Hleaf), root)
}
func Verify(str []byte, path []hash.Hash, leafchunk []byte, root hash.Hash) bool {
	return VerifyH(str, path, H(root, leafchunk), root)
}

