package main

import (
	"fmt"
	"flag"
	"math/rand"
	"encoding/hex"

	"time"

//	"merkletree"
	"merkletree/test_common"

	"signed_merkletree"
	"signed_merkletree/signed_merkletree_compliers"
)

//Add a `N` chunks and lists the tree leaves. `incp` is the probability of
// interest in a chunk.
func run_test(seed int64, n_min int32, n_max int32, N int, times int, subtimes int) {
	fmt.Println("Seed:", seed)
	r := rand.New(rand.NewSource(seed))
	
	gen := signed_merkletree.NewSignedMerkleProver()
	for i:= 0 ; i < N ; i++ {
		gen.AddChunk(test_common.Rand_chunk(r, n_min, n_max))
	}
	root := gen.Finish()  //Get the root hash.
	fmt.Println("Root:", hex.EncodeToString(root.Hash[:]))

	fmt.Println("---")

	// Set up signer.
	signer, pubkey := signed_merkletree_compliers.GenerateKey()

	j := 0
	for i:= 0 ; i < times ; i++ {
		// First part of challenge is a nonce.
		nonce := test_common.Rand_chunk(r, n_min, n_max)
		// Respond with root of the signed merkle tree.
		sigroot, smp:= gen.AddAllSigned(nonce, signer)
		for i2 := 0 ; i2 < subtimes ; i2++ {
			// (Nothing to check yet) Second part is randomly pick chunk.
			j := rand.Int63n(int64(N))
			
			// Create response:(the one that will be tested)
			// Regular and signature node.
			node, sig_node := gen.Getter.GetNode(j), smp.Getter.GetNode(j)
			chunk, sig := gen.Getter.GetChunk(j), smp.Getter.GetChunk(j)
			
			//Verify it.
			if !signed_merkletree.Verify(sig, chunk, nonce, pubkey, root.Hash, sigroot.Hash, node, sig_node) {
				fmt.Println("Didnt work")
			}
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
	var times int
	flag.IntVar(&times, "times", 256, "Number of times to challenge.")
	var subtimes int
	flag.IntVar(&subtimes, "subtimes", 4, "Number of indices per challenge.")
	
	flag.Parse()

	run_test(seed, int32(n_min), int32(n_max), N, times, subtimes)
}
