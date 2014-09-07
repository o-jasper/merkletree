package main

import (
	"fmt"
	"flag"
	"math/rand"

	"time"

	"merkle"
	"merkle/test_common"

	"signed_merkle"
	"signed_merkle/signed_merkle_pubkey"
	
	"crypto/sha256"
)

func main() {

//Get data.
	var seed, n_min, n_max, N, times, subtimes int64
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "Random seed for test.")
	flag.Int64Var(&n_min, "n_min", 1, "Minimum length of random chunk.")
	flag.Int64Var(&n_max, "n_max", 256, "Maximum length of random chunk.")
	flag.Int64Var(&N, "N", 80, "Number of chunks.")
	flag.Int64Var(&times, "times", 16, "Number of times to challenge.")
	flag.Int64Var(&subtimes, "subtimes", 8, "Number of indices per challenge.")
	var negative bool
	flag.BoolVar(&negative, "negative", false,
		"Wether to check for positives or try against false positives.")
	
	flag.Parse()

// Test Portion.
	fmt.Println("Denote info if it goes wrong.")
	fmt.Println("Seed:", seed)
	r := rand.New(rand.NewSource(seed))
	
	gen := signed_merkle.NewSignedMerkleProver(sha256.New(), false)
	for i:= int64(0) ; i < N ; i++ {
		gen.AddChunk(test_common.Rand_chunk(r, n_min, n_max))
	}
	root := gen.Finish()  //Get the root hash.
	fmt.Println("")
	fmt.Println("Root:", test_common.HashStr(root.Hash))

	fmt.Println("---")

	// Set up signer.
	signer, pubkey := signed_merkle_pubkey.GenerateKey()
	wrong_signer, _ := signed_merkle_pubkey.GenerateKey()

	for i:= int64(0) ; i < times ; i++ {
		// First part of challenge is a nonce.
		nonce := test_common.Rand_chunk(r, n_min, n_max)
		// Respond with root of the signed merkle tree.
		var sigroot *merkle.MerkleNode
		var smp *signed_merkle.SignedMerkleProver
		if !negative {
			sigroot, smp = gen.AddAllSigned(nonce, signer)
		} else {
			sigroot, smp = gen.AddAllSigned(nonce, wrong_signer)
		}

		for i2 := int64(0) ; i2 < subtimes ; i2++ {
			// (Nothing to check yet) Second part is randomly pick chunk.
			j := rand.Int63n(int64(N))
			
			// Create response:(the one that will be tested)
			// Regular and signature node.
			proof := gen.NewSignedMerkleProof_FromIndex(smp, j)
			//Verify it.
			if !negative {
				if r := proof.Verify(nonce, pubkey, root.Hash, sigroot.Hash); r != merkle.Correct {
					fmt.Println("Wrongly negative", r, ";", i, i2)
				}
			}	else if r := proof.Verify(nonce, pubkey, root.Hash, sigroot.Hash); r != merkle.WrongSig {
				if r == merkle.Correct { fmt.Println("False positive with wrong signer")
				} else { fmt.Println("Signature wasnt wrong with *just* that wrong.", r) }
			}
		}
	}
	fmt.Println("---")
	fmt.Println("No messages above implies success.")
	fmt.Println("times", times, "subtimes", subtimes, "N", N)
}
