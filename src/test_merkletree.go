package main

import (
	"fmt"
	"flag"
	"math/rand"
	"crypto/sha256"

	"time"

	"merkle"
	"merkle/test_common"

	"hash"
)

//Add a `N` chunks and lists the tree leaves. `incp` is the probability of
// interest in a chunk.
func main() {

//Read settings portion.
	var seed int64
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "Random seed for test.")
	var n_min, n_max, N int64
	flag.Int64Var(&n_min, "n_min", 1, "Minimum length of random chunk.")
	flag.Int64Var(&n_max, "n_max", 256, "Maximum length of random chunk.")
	flag.Int64Var(&N, "N", 256, "Number of chunks.")
	var incp float64
	flag.Float64Var(&incp, "incp", 0.3, "Probability of including to check.")
//	var bool 
//	flag.BoolVar(p *bool, name string, value bool, usage string)flag.
	flag.Parse()

	// Test portion.
	fmt.Println("Denote info if it goes wrong.")
	fmt.Println("Seed:", seed)
	r := rand.New(rand.NewSource(seed))

	gen := merkle.NewMerkleTreeGen(sha256.New(), false)  //Put chunks in.
	list := []*merkle.MerkleNode{}
	included := []bool{}
	for i:= int64(0) ; i < N ; i++ {
		chunk := test_common.Rand_chunk(r, n_min, n_max)
		include_this := (rand.Float64() <= incp)
		list = append(list, gen.AddChunk(chunk, include_this))
		included = append(included, include_this)
	}
	roothash := gen.Finish().Hash  //Get the root hash.
	fmt.Println("Root:", test_common.HashStr(roothash))

	fmt.Println("---")
//Reset random function, doing exact same to it.
	r = rand.New(rand.NewSource(seed))
	j, r2, ll := 0, rand.New(rand.NewSource(seed + 1)), 0
	for i := int64(0) ; i < N ; i++ {
		chunk := test_common.Rand_chunk(r, n_min, n_max)  // Recreates exactly as it was.
		root, valid := list[i].IsValid(-1)
		switch {
		case !valid:
			fmt.Println("Merkle tree not valid internally.")
		case !list[i].CorrespondsToChunk(chunk):
			fmt.Println("Chunk", i , "didnt check out.")
		case !root.CorrespondsToHash(roothash):
			fmt.Println("Not the correct top.", 
				test_common.HashStr(roothash), test_common.HashStr(root.Hash))
		default:
			r := list[i].Verify(roothash, chunk)
			if r != merkle.Correct {
				fmt.Println("Everything checked out but Verify didnt?", r)
			}
		}
		
		if included[i] {
			path := list[i].Path()
			ll = len(path)

			if !merkle.Verify(roothash, chunk, path) {
				fmt.Println(" - One of the Merkle Paths did not check out!")
				fmt.Println(test_common.HashStr(root.Hash))
			}
			j += 1
		}
		
		// Try if false positives might occur.
		if ll > 0 { // Use r2, other one needs to be the same as above!
			if r2.Int31() % 2 == 0 {
				chunk = test_common.Rand_chunk(r2, n_min, n_max)
			}
			path := []hash.Hash{}
			for i := 0 ; i < ll ; i++ {
				h := sha256.New() //10^-38 of accident, random gen will sooner collide.
				h.Write(test_common.Rand_bytes(r2, 16))
				path = append(path, h)
			}
			if merkle.Verify(roothash, chunk, path) {
				fmt.Println(" - False positive!")
			}
		}
	}
	fmt.Println("---")
	fmt.Println("No messages above implies success. Had", j)
}
