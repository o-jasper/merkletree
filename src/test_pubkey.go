package main

import (
	"time"
	"math/rand"
	"flag"
	"fmt"
	"signed_merkletree/signed_merkletree_pubkey"	
	"merkletree/test_common"

	"math/big"
)

func run_test(seed int64, n_min, n_max, times int32) {
	r := rand.New(rand.NewSource(seed))

	x := r.Int63()
	xb := big.NewInt(x).Bytes()
	if big.NewInt(0).SetBytes(xb).Cmp(big.NewInt(x)) != 0 {
		fmt.Println("Recreation failed ", x)

		fmt.Println(big.NewInt(x), big.NewInt(0).SetBytes(xb))
	}

	j := 0

	for i:= int32(0) ; i < times ; i++ {
		signer, pubkey := signed_merkletree_pubkey.GenerateKey()
		
		data := test_common.Rand_chunk(r, n_min, n_max)
		sig := signer.Sign(data)
		
		if !pubkey.VerifySignature(sig, data) {
			fmt.Print("*")
			j+= 1;
		}
	}
	fmt.Println("\nFailed", j, "of", times, "seed", seed)
}

func main() {
	var seed int64
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "Random seed for test.")
	var n_min int64
	flag.Int64Var(&n_min, "n_min", 1, "Minimum length of random chunk.")
	var n_max int64
	flag.Int64Var(&n_max, "n_max", 10, "Maximum length of random chunk.")
	var times int64
	flag.Int64Var(&times, "times", 80, "Number of times to challenge.")

	run_test(seed, int32(n_min), int32(n_max), int32(times))
}
