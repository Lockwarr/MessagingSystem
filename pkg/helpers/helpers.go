package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
)

// GeneratePublicKey - generates public key
func GeneratePublicKey() string {
	key, err := rsa.GenerateKey(rand.Reader, 128)
	if err != nil {
		fmt.Println(err)
	}
	// Generate public key Main address
	PublicKey := &key.PublicKey
	//The public key is generated from the private key
	pkixPublicKey, err := x509.MarshalPKIXPublicKey(PublicKey)
	return string(pkixPublicKey)
}
