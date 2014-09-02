package main

import (
	"fmt"
	"math/rand"
	"trie_easy"

	"merkletree/test_common"
	"encoding/hex"

	"time"
)

func main() {

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	node := trie_easy.NewTrie(nil)
	compare := map[string]int{}
	
	n_min, n_max := int32(3), int32(8)
	
	fmt.Println("= Put stuff in twice.")
	for i := 0 ; i < 10 ; i++ {
		k := test_common.Rand_chunk(r, n_min, n_max)
		fmt.Println(i, hex.EncodeToString(k))
		node.Set(k, i, trie_easy.Creator16{})
		compare[string(k)] = i
	}
	fmt.Println("= Printing")
	node.MapAll(nil, func(data interface{}, k []byte, v interface{}) bool{
		fmt.Println(v, hex.EncodeToString(k))
		return false
	})
	fmt.Println("= And check equality.")
	for k,v := range compare {
		val := node.Get([]byte(k), 0)
		if val != v {
			fmt.Println("Mismatch on:", hex.EncodeToString([]byte(k)),":", v, "vs", val)
		}
	}
	//fmt.Print("= And lack of presence") //TODO
}
