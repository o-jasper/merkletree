package main

import (
	"fmt"
	"math/rand"
	"trie_easy"
	"merkletree/test_common"
)

func prt16(n *trie_easy.TrieNode16) {
	fmt.Print("(")
	for i, el := range n.Sub {
		if data, ok := el.Actual.(*trie_easy.TrieNode16) ; ok && data.Data!=nil {
			fmt.Print(i, ",") 
		}
		prt(&el)
	}
	fmt.Print(")")
}

func prt(node* trie_easy.TrieNode) {
	c, okc := node.Actual.(*trie_easy.TrieNode16)
	if okc { 
		prt16(c) 
	}
}

func main() {

	seed := int64(243525623)
	r := rand.New(rand.NewSource(seed))
	node := trie_easy.NewTrieNode(nil)
	compare := map[string][]byte{}
	
	n_min, n_max := int32(3), int32(8)
	
	fmt.Println("= Put stuff in twice.")
	for i := 0 ; i < 100 ; i++ {
		set,to := test_common.Rand_chunk(r, n_min, n_max), test_common.Rand_chunk(r, n_min, n_max)
		node.Set(set, to)
		compare[string(set)] = to
	}
	fmt.Println("= And check equality.")
	for k,v := range compare {
		got := node.Get([]byte(k), 0)
		val, ok := got.([]byte)
		if ok && string(val) != string(v) {
			fmt.Println("Mismatch on:", k,":", v, "vs", val)
		} else if !ok {
			fmt.Println("Couldnt convert back", val, v)
		}
	}
}
