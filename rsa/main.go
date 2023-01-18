package main

import (
	"fmt"
	"crypto"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey

	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(), // hashing function - SHA256 algorithm
		rand.Reader,
		&publicKey,
		[]byte("super secret message"),
		nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("encrypted bytes: ", encryptedBytes)

	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})

	fmt.Println("decrypted bytes: ", string(decryptedBytes))

	msg := []byte("verifiable message")

	msgHash := sha256.New()
	_, err1 := msgHash.Write(msg)
	if err1 != nil {
		panic(err1)
	}
	msgHashSum := msgHash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}

	err = rsa.VerifyPSS(&publicKey, crypto.SHA256, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return
	}

	fmt.Println("signature verified")
}
