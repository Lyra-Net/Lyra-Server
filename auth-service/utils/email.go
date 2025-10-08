package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
)

type CryptoManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

var Crypto *CryptoManager

func InitCrypto() {
	priv, err := LoadPrivateKey("/keys/private.pem")
	if err != nil {
		log.Fatal("failed to load private key:", err)
	}
	pub, err := LoadPublicKey("/keys/public.pem")
	if err != nil {
		log.Fatal("failed to load public key:", err)
	}

	Crypto = &CryptoManager{
		privateKey: priv,
		publicKey:  pub,
	}
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid private key PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}

	return rsaKey, nil
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid public key PEM")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return rsaPub, nil
}

func EncryptEmailRSA(email string) (string, error) {
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, Crypto.publicKey, []byte(email))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func DecryptEmailRSA(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, Crypto.privateKey, data)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func HashEmail(email string) string {
	h := sha256.Sum256([]byte(email))
	return hex.EncodeToString(h[:])
}
