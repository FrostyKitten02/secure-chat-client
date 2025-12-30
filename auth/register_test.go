package auth

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func TestEncDecPrivKey(t *testing.T) {
	curve := ecdh.X25519()

	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
		return
	}

	password := "password123123123"
	encPrivKey, privKeyEncErr := encryptPrivateKey(password, privateKey.Bytes())
	if privKeyEncErr != nil {
		t.Fatal(privKeyEncErr)
		return
	}

	decPrivKey, decPrivKeyErr := decryptPrivateKey(password, encPrivKey)
	if decPrivKeyErr != nil {
		t.Fatal(decPrivKeyErr)
		return
	}
	noChange := bytes.Equal(privateKey.Bytes(), decPrivKey)
	if !noChange {
		t.Fatal("Private keys don't match")
		return
	}

	base64EncKey := base64.StdEncoding.EncodeToString(encPrivKey)
	originalEncKey, decodeErr := base64.StdEncoding.DecodeString(base64EncKey)
	if decodeErr != nil {
		t.Fatal(decodeErr)
		return
	}

	base64DecKey, decrpytBase64Err := decryptPrivateKey(password, originalEncKey)
	if decrpytBase64Err != nil {
		t.Fatal(decrpytBase64Err)
		return
	}

	base64NoChange := bytes.Equal(base64DecKey, privateKey.Bytes())
	if !base64NoChange {
		t.Fatal("Decrypted private keys don't match")
		return
	}
}
