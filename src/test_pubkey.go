package main

import (
	"time"
	"math/rand"
	"flag"
	"fmt"
	"signed_merkle/signed_merkle_pubkey"	
	"merkle/test_common"

	"math/big"
)

func main() {

// Get data portion.
	var seed, n_min, n_max, times int64
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "Random seed for test.")
	flag.Int64Var(&n_min, "n_min", 1, "Minimum length of random chunk.")
	flag.Int64Var(&n_max, "n_max", 10, "Maximum length of random chunk.")
	flag.Int64Var(&times, "times", 80, "Number of times to challenge.")

// Running portion.
	r := rand.New(rand.NewSource(seed))

	x := r.Int63()
	xb := big.NewInt(x).Bytes()
	if big.NewInt(0).SetBytes(xb).Cmp(big.NewInt(x)) != 0 {
		fmt.Println("Recreation failed ", x)

		fmt.Println(big.NewInt(x), big.NewInt(0).SetBytes(xb))
	}

	j := 0

	for i:= int64(0) ; i < times ; i++ {
		signer, pubkey := signed_merkle_pubkey.GenerateKey()
		wrong_signer, _ := signed_merkle_pubkey.GenerateKey()
		
		data := test_common.Rand_chunk(r, n_min, n_max)
		sig := signer.Sign(data)
		wrong_sig := wrong_signer.Sign(data)
		
		if !pubkey.VerifySignature(sig, data) {
			fmt.Print("This one should have been correct.")
			j+= 1;
		}
		if pubkey.VerifySignature(wrong_sig, data) {
			fmt.Print("This one should have been incorrect.")
		}
		// TODO check against false positive.
	}
	fmt.Println("\nFailed", j, "of", times, "seed", seed)
}
