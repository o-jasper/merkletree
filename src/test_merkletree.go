package main

import (
	"fmt"
	"flag"
	"math/rand"
	"crypto/sha256"

//	"time"

	"merkle"
	"merkle/test_common"

	"hash_extra"
)

//Add a `N` chunks and lists the tree leaves. `incp` is the probability of
// interest in a chunk.
func main() {

//Read settings portion.
	var seed, n_min, n_max, N int64
	flag.Int64Var(&seed, "seed", 1/*time.Now().UnixNano()*/, "Random seed for test.")
	flag.Int64Var(&n_min, "n_min", 1, "Minimum length of random chunk.")
	flag.Int64Var(&n_max, "n_max", 256, "Maximum length of random chunk.")
	flag.Int64Var(&N, "N", 256, "Number of chunks.")
	var incp float64
	flag.Float64Var(&incp, "incp", 1, "Probability of including to check.")
	var with_index bool 
	flag.BoolVar(&with_index, "with_index", true, "Wether to have an index in each chunk")
	flag.Parse()

	// Test portion.
	fmt.Println("Denote info if it goes wrong.")
	fmt.Println("Seed:", seed)
	r := rand.New(rand.NewSource(seed))

	hasher := hash_extra.Hasher{sha256.New()}
	gen := merkle.NewMerkleTreeGen(hasher, with_index)  //Put chunks in.
	list := []*merkle.MerkleNode{}
	included := []bool{}
	for i:= int64(0) ; i < N ; i++ {
		chunk := test_common.Rand_chunk(r, n_min, n_max)
		include_this := (rand.Float64() <= incp)
		list = append(list, gen.Add(chunk, include_this))
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
		root, valid := list[i].IsValid(hasher, -1)
		switch {
		case !valid:
//			fmt.Println(test_common.HashStr(root.Left.Hash), test_common.HashStr(root.Right.Hash))
//			fmt.Println(test_common.HashStr(root.Hash), test_common.HashStr(hasher.H_U2(root.Right.Hash, root.Left.Hash)))
			fmt.Println("Merkle tree not valid internally.")
			
		case with_index && !list[i].CorrespondsWithIndex(hasher, uint64(i), chunk) || !with_index && !list[i].Corresponds(hasher, chunk):
			fmt.Println("Chunk", i , "didnt check out.")
		case included[i] && !root.CorrespondsH(roothash):
			fmt.Println("Not the correct top.", root.Up)
//				test_common.HashStr(roothash), test_common.HashStr(root.Hash))
		case included[i]:
			var r int8
			if with_index { r = list[i].VerifyWithIndex(hasher, roothash, uint64(i), chunk)
			}	else { r = list[i].Verify(hasher, roothash, chunk) }
			if r != merkle.Correct {
				fmt.Println("Everything checked out but Verify didnt?", r)
			}
		}
		
		if included[i] {
			path, success := list[i].Path(), false
			ll = len(path)
			if with_index { 
				success = hasher.MerkleVerifyWithIndex(roothash, uint64(i), chunk, path)
			} else { success = hasher.MerkleVerify(roothash, chunk, path) }
			
			if !success {
				fmt.Println(" - One of the Merkle Paths did not check out!", ll)
				fmt.Println(
					test_common.HashStr(hasher.MerkleExpectedRoot(hasher.HwI(uint64(i), chunk), path)))
			}
			j += 1
		}
		
		// Try if false positives might occur.
		if ll > 0 { // Use r2, other one needs to be the same as above!
			if r2.Int31() % 2 == 0 {
				chunk = test_common.Rand_chunk(r2, n_min, n_max)
			}
			path := []hash_extra.HashResult{}
			for i := 0 ; i < ll ; i++ {
				//10^-38 of accident, random gen will sooner collide.
				path = append(path, hasher.H(test_common.Rand_bytes(r2, 16)))
			}
			if hasher.MerkleVerify(roothash, chunk, path) {
				fmt.Println(" - False positive!")
			}
		}
	}
	fmt.Println("---")
	fmt.Println("No messages above implies success. Had", j, " (TODO number of tests.. i reckon)")
}
