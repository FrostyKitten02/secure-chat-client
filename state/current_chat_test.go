package state

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"testing"
)

func Test_shared_secret(t *testing.T) {
	pub1, priv1 := createKeyPair(t)
	pub2, priv2 := createKeyPair(t)

	shared2, s2err := createChatKey(pub1.Bytes(), priv2.Bytes())
	if s2err != nil {
		t.Errorf("Failed to create shared secret: %s", s2err)
		return
	}

	shared1, s1err := createChatKey(pub2.Bytes(), priv1.Bytes())
	if s1err != nil {
		t.Errorf("Failed to create shared secret: %s", s1err)
		return
	}

	if !bytes.Equal(shared1, shared2) {
		t.Errorf("Shared secrets don't match")
	}

}

func createKeyPair(t *testing.T) (privateKey *ecdh.PrivateKey, publicKey *ecdh.PublicKey) {
	curve := ecdh.X25519()
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		t.Error("failed to generate keys")
		return nil, nil
	}

	pubKey := privateKey.PublicKey()
	return privateKey, pubKey
}
