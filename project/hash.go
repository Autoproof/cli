package project

import (
	"encoding/hex"
	"hash"

	"golang.org/x/crypto/sha3"
)

type keccakHash struct {
	hash.Hash
}

func newKeccakHash() *keccakHash {
	return &keccakHash{Hash: sha3.NewLegacyKeccak256()}
}

func (k *keccakHash) String() string {
	return "keccak256:" + hex.EncodeToString(k.Sum(nil))
}
