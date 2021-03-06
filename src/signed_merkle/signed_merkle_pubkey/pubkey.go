//Note: really this is basically more general than just for this thing...
// Some of go's standard lib seems to be a bit too limited.
package signed_merkle_pubkey

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"

	"math/big"
)

// Signature stuff.
type EcdsaSigner struct {
	Priv *ecdsa.PrivateKey
}

func (signer EcdsaSigner) Sign(input []byte) []byte {
	r, s, _ := ecdsa.Sign(rand.Reader, signer.Priv, input)
	rd := r.Bytes()
	sd := s.Bytes()
	
	lens := []byte{byte(len(rd))}
	return append(append(lens, rd...), sd...)
}

type EcdsaPubkey struct {
	Pub *ecdsa.PublicKey
}

func (pubkey EcdsaPubkey) VerifySignature(sig []byte, data []byte) bool {
	lr := sig[0]
	r := big.NewInt(0).SetBytes(sig[1 : 1 + lr])
	s := big.NewInt(0).SetBytes(sig[1 + lr :])
	return ecdsa.Verify(pubkey.Pub, data, r, s)
}

//WARNING ensure you got a good random source AND elliptic thingy!
func GenerateKey() (EcdsaSigner, EcdsaPubkey) {
	//NOTE this one is random, pick good ones!
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return EcdsaSigner{Priv : priv}, EcdsaPubkey{Pub : &priv.PublicKey}
}
