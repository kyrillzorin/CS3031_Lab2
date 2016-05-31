package main

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Load RSA private key from pem file
func privateKeyFromFile() (*rsa.PrivateKey, error) {
	var err error
	var pemData []byte
	var block *pem.Block
	var privateKey *rsa.PrivateKey
	if pemData, err = ioutil.ReadFile("./priv.pem"); err != nil {
		err = fmt.Errorf("Error reading pem file: %s", err)
		return nil, err
	}
	if block, _ = pem.Decode(pemData); block == nil || block.Type != "RSA PRIVATE KEY" {
		err = errors.New("No valid PEM data found")
		return nil, err
	}
	if privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		err := fmt.Errorf("Private key can't be decoded: %s", err)
		return nil, err
	}
	return privateKey, nil
}

// Generate a new RSA private key and save to pem file
func generatePrivateKey() error {
	var privateKey *rsa.PrivateKey
	var err error
	if privateKey, err = rsa.GenerateKey(rand.Reader, 1024); err != nil {
		return err
	}
	privateKeyFile, err := os.Create("./priv.pem")
	defer privateKeyFile.Close()
	if err != nil {
		return err
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}
	return nil
}

// Get RSA private key from file or generate new one if necessary
func getPrivateKey() (*rsa.PrivateKey, error) {
	if _, err := os.Stat("./priv.pem"); os.IsNotExist(err) {
		err = generatePrivateKey()
		if err != nil {
			return nil, err
		}
	}
	privateKey, err := privateKeyFromFile()
	return privateKey, err
}

// Encrypt data using RSA public key
func encrypt(public_key *rsa.PublicKey, plain_text []byte) ([]byte, error) {
	var label, encrypted []byte
	var err error
	if encrypted, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, public_key, plain_text, label); err != nil {
		return nil, err
	}
	return encrypted, nil
}

// Decrypt data using RSA private key
func decrypt(private_key *rsa.PrivateKey, encrypted []byte) ([]byte, error) {
	var label, decrypted []byte
	var err error
	if decrypted, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, private_key, encrypted, label); err != nil {
		return nil, err
	}
	return decrypted, nil
}

// Sign message using RSA private key
func sign(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	hasher := crypto.SHA256.New()
	hasher.Write(message)
	hashed := hasher.Sum(nil)
	var opts rsa.PSSOptions
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hashed, &opts)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Encrypt data using AES key
func encryptAES(key, data []byte) ([]byte, error) {
	var err error
	var block cipher.Block
	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}
	encryptedData := make([]byte, aes.BlockSize+len(data))
	initVector := encryptedData[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, initVector); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, initVector)
	stream.XORKeyStream(encryptedData[aes.BlockSize:], data)
	return encryptedData, nil
}

// Decrypt data using AES key
func decryptAES(key, encryptedData []byte) ([]byte, error) {
	var err error
	var block cipher.Block
	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	if len(encryptedData) < aes.BlockSize {
		err = errors.New("encryptedData too short")
		return nil, err
	}
	initVector := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, initVector)
	stream.XORKeyStream(encryptedData, encryptedData)
	return encryptedData, nil
}

// Generate new AES (256) key
func generateAESKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}
