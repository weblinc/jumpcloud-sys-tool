package jc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
)

type ClientPrivateKey struct {
	*rsa.PrivateKey
}

type ClientPrivateKeySigner interface {
	Sign(data []byte) ([]byte, error)
	SignatureForRequest(string, string, string) (string, error)
}

// Loads a RSA private key from file, base64 decrypts the PEM, parses and returns an rsa.PrivateKey struct
func LoadClientPrivateKeyFromFile(path string) (pk ClientPrivateKey, err error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	block, _ := pem.Decode(contents)
	if block == nil {
		err = errors.New("No key found")
		return
	}

	var rawKey interface{}

	// Parse the private key from PEM block
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal("PK Parse error:", err)
		}
		rawKey = rsa
	default:
		log.Fatal("unsupported priv key type", block.Type)
	}

	switch t := rawKey.(type) {
	case *rsa.PrivateKey:
		pk = ClientPrivateKey{t}
	default:
		log.Println("No RSA priv key :(")
	}

	return
}

// Signs a request signature with the current private key
func (pk *ClientPrivateKey) SignatureForRequest(time, httpMethod, httpUrl string) (signedBase64 string, err error) {
	msg := httpMethod + " " + httpUrl + " HTTP/1.1\ndate: " + time
	signedMsg, err := pk.Sign([]byte(msg))
	if err != nil {
		return
	}

	signedBase64 = base64.StdEncoding.EncodeToString(signedMsg)
	return
}

// SHA256 hashes an array of bytes and signs it with the current private key
func (pk *ClientPrivateKey) Sign(data []byte) (signed []byte, err error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, pk.PrivateKey, crypto.SHA256, d)
}
