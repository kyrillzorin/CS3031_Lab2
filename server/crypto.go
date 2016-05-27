package main

import (
	"crypto"
	"crypto/rsa"
)

func verify(publicKey *rsa.PublicKey, message []byte, signature []byte) bool {
	hasher := crypto.SHA256.New()
	hasher.Write(message)
	hashed := hasher.Sum(nil)
	err := rsa.VerifyPSS(publicKey, crypto.SHA256, hashed, signature, nil)
	if err != nil {
		return false
	}
	return true
}
