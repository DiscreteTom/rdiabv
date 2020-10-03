package main

import (
	"crypto/rsa"
	"io"
	"math/big"
)

// RawRsa will store an rsa key pair.
type RawRsa struct {
	sk *rsa.PrivateKey
}

// NewRawRsa will generate a key pair.
// If you don't know which random to use, use rand.Reader.
func NewRawRsa(random io.Reader, bits int) (rr *RawRsa) {
	rr = &RawRsa{}
	var sk, err = rsa.GenerateKey(random, bits)
	if err != nil {
		panic("Private key generation failed.")
	}
	rr.sk = sk
	return
}

// Encrypt will encrypt the given secretMsg.
func (rr *RawRsa) Encrypt(secretMsg *big.Int) (ciphertext *big.Int) {
	var sk = rr.sk
	ciphertext = &big.Int{}
	ciphertext.Exp(secretMsg, new(big.Int).SetInt64(int64(sk.E)), sk.N)
	return
}

// Decrypt will decrypt the given ciphertext.
func (rr *RawRsa) Decrypt(ciphertext *big.Int) (secretMsg *big.Int) {
	var sk = rr.sk
	secretMsg = &big.Int{}
	secretMsg.Exp(ciphertext, sk.D, sk.N)
	return
}
