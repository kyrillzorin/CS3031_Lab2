package main

import (
	"crypto"
	"crypto/rsa"
)

// Verify an RSA signed message
func verify(publicKey *rsa.PublicKey, message []byte, signature []byte) bool {
	hasher := crypto.SHA256.New()
	hasher.Write(message)
	hashed := hasher.Sum(nil)
	var opts rsa.PSSOptions
	err := rsa.VerifyPSS(publicKey, crypto.SHA256, hashed, signature, &opts)
	if err != nil {
		return false
	}
	return true
}
