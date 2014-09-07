package main

import (
	"fmt"
	"math/rand"
	"trie_easy"

	"merkle/test_common"
	"encoding/hex"

	"time"
)

func main() {

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	node := trie_easy.NewTrie(nil)
	compare := map[string]int{}
	
	n_min, n_max, N := int32(3), int32(8), 10
	
	fmt.Println("= Put stuff in twice.")
	for i := 0 ; i < N ; i++ {
		chunk := test_common.Rand_chunk(r, n_min, n_max)
		fmt.Println(i, hex.EncodeToString(chunk))
		node.Set(chunk, i, trie_easy.Creator16{})
		compare[string(chunk)] = i
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
	fmt.Println("= Try against false positives.")
	for i := 0 ; i < N ; i++ {
		key := test_common.Rand_chunk(r, n_min, n_max)
		v, got := compare[string(key)]
		val := node.Get([]byte(key), 0)
		if got { // Accidentally made one that exists.
			if val != v {
				fmt.Println("(during rand chunks)Mismatch on:", 
					hex.EncodeToString([]byte(key)),":", v, "vs", val)
			}
		} else if val != nil {
			fmt.Println("Got something extra:", val)
		}
	}
	//fmt.Print("= And lack of presence") //TODO
}
