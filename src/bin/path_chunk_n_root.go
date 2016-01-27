package main

import (
	"fmt"
	"flag"
	"math/rand"
	"merkle"
	"encoding/hex"
	
	common "test/common"
	"hash_extra"
	"crypto/sha256"
)

//Returns:
// * chunk itself
// * chunk hash
// * path
// * root hash

func gen_data(seed int64, n_min int64, n_max int64, N int, i int) {
	r := rand.New(rand.NewSource(seed))

	gen := merkle.NewMerkleTreeGen(hash_extra.Hasher{sha256.New()}, true)  //Put chunks in.
	var node *merkle.MerkleNode
	node = nil
	for j:= 0 ; j < N ; j++ {
		chunk := common.Rand_chunk(r, n_min, n_max)
		if j == i {
			fmt.Println(hex.EncodeToString(chunk))  //Print the chunk itself.			
			node = gen.Add(chunk, true)
		} else { 
			gen.Add(chunk, false)
		}
	}
	fmt.Println(hex.EncodeToString(node.Hash[:])) //Print the hash of the chunk.
	root_hash := gen.Finish().Hash[:]
	
	path := node.Path() //Print the path.
	for j := range path {
		fmt.Print(hex.EncodeToString(path[j][:]))
	}
	fmt.Println()
	//Print the root.
	fmt.Println(hex.EncodeToString(root_hash[:]))
}

func main() {
	var seed int64
	flag.Int64Var(&seed, "seed", rand.Int63(), "Random seed for test.")
	var n_min int64
	flag.Int64Var(&n_min, "n_min", 1, "Minimum length of random chunk.")
	var n_max int64
	flag.Int64Var(&n_max, "n_max", 256, "Maximum length of random chunk.")
	var N int
	flag.IntVar(&N, "N", 256, "Number of chunks.")
	r := rand.New(rand.NewSource(seed))
	var i int
	flag.IntVar(&i, "i", int(common.Rand_range(r, 0, int64(N-1))), "Which chunk to get.")
	
	flag.Parse()

	gen_data(r.Int63(), int64(n_min), int64(n_max), N, i)
}

