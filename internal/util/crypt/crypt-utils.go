// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// GetKeys generates and returns a private and public key pair
func GetKeys() (string, string, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		return "", "", err
	}

	err = privateKey.Validate()
	if err != nil {
		return "", "", err
	}

	// Get der format. priv_der []byte
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// PEM Block
	privateKeyPEMBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDER,
	}

	// Public Key generation
	publicKey := privateKey.PublicKey
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return "", "", err
	}

	// PEM Block
	publicKeyPEMBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDER,
	}

	// Resultant keys are in PEM format.
	return string(pem.EncodeToMemory(&privateKeyPEMBlock)), string(pem.EncodeToMemory(&publicKeyPEMBlock)), err
}

// GetHash returns a unique hash for a given string
func GetHash(text string) string {
	return fmt.Sprintf("%x", sha512.Sum512_256([]byte(text)))
}
