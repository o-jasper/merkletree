package main

import (
	"fmt"
	"flag"
	"math/rand"
	"trie_easy"

	"merkle/test_common"
	"merkle/merkle_common"
	"encoding/hex"

	"time"

	"trie_merkle"
	"crypto/sha256"

	"hash"
	"hash_extra"
)

func main() {

	var seed, n_min, n_max, N int64
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "Random seed for test.")
	flag.Int64Var(&n_min, "n_min", 0, "Minimum length of random chunk.")
	flag.Int64Var(&n_max, "n_max", 64, "Maximum length of random chunk.")
	flag.Int64Var(&N, "N", 80, "Number of entries.")
	var print, merkle bool
	flag.BoolVar(&print, "print", false, "Whether to print everything")
	flag.BoolVar(&merkle, "merkle", true, "Whether to do the merkle part")

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
	var root hash.Hash
	if merkle { 
		fmt.Println("= Check equality, and verification through path.")
		root = trie_merkle.Hash(&node, sha256.New())
		fmt.Println("= Got root", test_common.HashStr(root))
	} else {
		root = nil
		fmt.Println("= And check equality.")
	}
	for k,v := range compare {
		val := node.Get([]byte(k), 0)
		if val != v {
			fmt.Println("Mismatch on:", hex.EncodeToString([]byte(k)),":", v, "vs", val)
		}
		if merkle { // Run the test.
			// TODO
			str := []byte(k)
			_, _, path := trie_merkle.HashPath(&node, str, sha256.New())
			h := hash_extra.H(sha256.New(), merkle_common.GetBytes(v))
			if !trie_merkle.VerifyH(str, path, h, root) {
				fmt.Println("Verification didnt work.")
			}
			// TODO also, try against false positives.
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
}
