//Note: really this is basically more general than just for this thing...
// (why didnt go come like this?.. Unclear to me)
package signed_merkletree_pubkey

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
	r := big.NewInt(0).SetBytes(sig[2 : 2 + lr])
	s := big.NewInt(0).SetBytes(sig[2 + lr :])
	return ecdsa.Verify(pubkey.Pub, data, r, s)
}

func GenerateKey() (EcdsaSigner, EcdsaPubkey) {
	//NOTE this one is random, pick good ones!
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return EcdsaSigner{Priv : priv}, EcdsaPubkey{Pub : &priv.PublicKey}
}
