package main

import (
	"fmt"
	"flag"
	"math/rand"
//	"crypto/sha256"
	"encoding/hex"

	"time"

	"merkletree"
	"merkletree/test_common"
)

//Add a `N` chunks and lists the tree leaves. `incp` is the probability of
// interest in a chunk.
func run_test(seed int64, n_min int32, n_max int32, N int, incp float64) {
	fmt.Println("Seed:", seed)
	r := rand.New(rand.NewSource(seed))

	gen := merkletree.NewMerkleTreeGen()  //Put chunks in.
	list := []*merkletree.MerkleNode{}
	included := []bool{}
	for i:= 0 ; i < N ; i++ {
		chunk := test_common.Rand_chunk(r, n_min, n_max)
		include_this := (rand.Float64() <= incp)
		list = append(list, gen.AddChunk(chunk, include_this))
		included = append(included, include_this)
	}
	roothash := gen.Finish().Hash  //Get the root hash.
	fmt.Println("Root:", hex.EncodeToString(roothash[:]))

	fmt.Println("---")
//Reset random function, doing exact same to it.
	r = rand.New(rand.NewSource(seed))
	j := 0
	for i:= 0 ; i < N ; i++ {
		chunk := test_common.Rand_chunk(r, n_min, n_max)
		root, valid := list[i].IsValid(-1)
		switch {
		case !valid:                             fmt.Println("Merkle tree not valid internally.")
		case !list[i].CorrespondsToChunk(chunk): fmt.Println("Chunk", i , "didnt check out.")
		case !root.CorrespondsToHash(roothash):
			fmt.Println("Not the correct top.", 
				hex.EncodeToString(roothash[:]), hex.EncodeToString(root.Hash[:]))
		default:
			if r := list[i].Verify(roothash, merkletree.H(chunk)); r != merkletree.Correct {
				fmt.Println("Everything checked out but Verify didnt?", r)
			}
		}
		
		if included[i] {
			path := list[i].Path()
// For if you want to print it.
//		root := merkletree.ExpectedRoot(merkletree.H(chunk), path)
//		fmt.Println(hex.EncodeToString(root[:]))
			
			if !merkletree.Verify(roothash, chunk, path) {
				fmt.Println(" - One of the Merkle Paths did not check out!")
			}
			j += 1
		}
	}
	fmt.Println("---")
	fmt.Println("No messages above implies success. Had", j)
}

func main() {
	var seed int64
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "Random seed for test.")
	var n_min int64
	flag.Int64Var(&n_min, "n_min", 1, "Minimum length of random chunk.")
	var n_max int64
	flag.Int64Var(&n_max, "n_max", 256, "Maximum length of random chunk.")
	var N int
	flag.IntVar(&N, "N", 256, "Number of chunks.")
	var incp float64
	flag.Float64Var(&incp, "incp", 0.3, "Probability of including to check.")
	
	flag.Parse()

	run_test(seed, int32(n_min), int32(n_max), N, incp)
}
