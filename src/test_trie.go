package main

import (
	"fmt"
	"flag"
	"math/rand"
	"trie_easy"

	"merkle/test_common"
	"encoding/hex"

	"time"

	"trie_merkle"
	"crypto/sha256"
)

func main() {

	var seed, n_min, n_max, N int64
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "Random seed for test.")
	flag.Int64Var(&n_min, "n_min", 0, "Minimum length of random chunk.")
	flag.Int64Var(&n_max, "n_max", 64, "Maximum length of random chunk.")
	flag.Int64Var(&N, "N", 80, "Number of entries.")
	var print bool
	flag.BoolVar(&print, "print", false, "Whether to print everything")

	node := trie_easy.NewTrie(nil)
	compare := map[string]int64{}

	r := rand.New(rand.NewSource(seed))
	
	fmt.Println("= Put stuff in twice.")
	for i := int64(0) ; i < N ; i++ {
		chunk := test_common.Rand_chunk(r, n_min, n_max)
		if print { fmt.Println(i, hex.EncodeToString(chunk)) }
		node.Set(chunk, i, trie_easy.Creator16{})
		compare[string(chunk)] = i
	}
	if print {
		fmt.Println("= Printing")
		node.MapAll(nil, func(data interface{}, k []byte, v interface{}) bool{
			fmt.Println(v, hex.EncodeToString(k))
			return false
		})
	}
	fmt.Println("= And check equality.")
	for k,v := range compare {
		val := node.Get([]byte(k), 0)
		if val != v {
			fmt.Println("Mismatch on:", hex.EncodeToString([]byte(k)),":", v, "vs", val)
		}
	}
	fmt.Println("= Try against false positives.")
	for i := int64(0) ; i < N ; i++ {
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

	fmt.Println(test_common.HashStr(trie_merkle.Hash(&node, sha256.New())))
}
