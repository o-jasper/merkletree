package data_n_privkey

//WARNING about using this for Proof of Custody - style stuff.
// Some pubkey signing algos might reduce to a checksum being signed, so people
// can cheat by sending checksums to fake having the data.
//(the purpose of showing you have both data and (solely)privkey fails)

import (
	"merkletree"
//	"crypto"
	"crypto/sha256"
)

type DataNPrivKey struct {
	Signer interface{}
	Sigs merkletree.MerkleTreeGen
	Data merkletree.MerkleTreeGen
}

type DataNPrivKeyInterest struct {
	Chunk     []byte
	ChunkSig  []byte
	Sigs  *merkletree.MerkleNode
	Data  *merkletree.MerkleNode
}

func (dpk *DataNPrivKey) AddChunk(chunk []byte, interest bool) *DataNPrivKeyInterestRoots {
	sig := dpk.Signer.Sign(chunk)
	H := merkletree.H(chunk)
	
	return &DataNPrivKeyInterest{
		Chunk    : chunk,
		ChunkSig : sig,
		Sigs  : dpk.Sigs.AddChunk(sig, interest),
		Data  : dpk.Data.AddChunkH(H, interest)
	}
}

func (dpk *DataNPrivKey) Finish() DataNPrivKeyRoots {
	return DataNPrivKeyRoots{Sigs:dpk.Sigs.Finish(), Data:dpk.Sigs.Finish()}
}

func (interest *DataNPrivKeyInterest) IsValid(recurse int32) *DataNPrivKeyInterest, bool {
	s, svalid := interest.Sigs.IsValid(recurse)
	d, dvalid := interest.Data.IsValid(recurse)  // NOTE crosscheck left/right paths sequence?
	return &DataNPrivKeyInterest{Sigs:s, Data:d}, svalid && dvalid
}

func (interest *DataNPrivKeyInterest) CorrespondsToPubkey(pubkey interface{}) bool {
	return pubkey.Verify(interest.ChunkSig, interest.Chunk)
}

func (interest *DataNPrivKeyInterest) Verify(Hts [sha256.Size]byte, Htd [sha256.Size]byte, pubkey interface) bool {
	root, internal := interest.IsValid(-1)
	return (internal && 
		interest.CorrespondsToPubkey(pubkey) && 
		root.Sigs.CorrespondsToHash(Hts) && 
		root.Data.CorrespondsToHash(Htd) &&
		interest.Sig.CorrespondsToChunk(Interest.ChunkSig) &&
		interest.Data.CorrespondsToChunk(Interest.Chunk))
}
