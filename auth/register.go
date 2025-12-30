package auth

import (
	"crypto"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`

	PubKey     string `json:"pubKey"`
	EncPrivKey string `json:"encPrivKey"`
}

func Register(username, email, password string) error {
	curve := ecdh.X25519()

	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return errors.New("failed to generate keys")
	}

	encryptedPrivKey, privKeyErr := encryptPrivateKey(password, privateKey.Bytes())
	if privKeyErr != nil {
		return errors.New("failed to encrypt private key")
	}
	publicKey := privateKey.PublicKey()

	req := &RegisterRequest{
		Email:    email,
		Username: username,
		Password: password,

		PubKey:     string(publicKey.Bytes()),
		EncPrivKey: base64.StdEncoding.EncodeToString(encryptedPrivKey),
	}
	//TODO: make api call!
	if req != nil {

	}

	return nil
}

func encryptPrivateKey(password string, privKey []byte) ([]byte, error) {
	chaChaKey, nonce, chaChaErr := getChaCha20Poly1305Params(password)
	if chaChaErr != nil {
		return nil, chaChaErr
	}

	aead, err := chacha20poly1305.New(chaChaKey)
	if err != nil {
		return nil, err
	}

	return aead.Seal(nil, nonce, privKey, nil), nil
}

func decryptPrivateKey(password string, encPrivKey []byte) ([]byte, error) {
	chaChaKey, nonce, chaChaErr := getChaCha20Poly1305Params(password)
	if chaChaErr != nil {
		return nil, chaChaErr
	}

	aead, err := chacha20poly1305.New(chaChaKey)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, nonce, encPrivKey, nil)
}

func getChaCha20Poly1305Params(password string) ([]byte, []byte, error) {
	secret, keyGenErr := hkdf.Key(crypto.SHA256.New, []byte(password), nil, "", chacha20poly1305.KeySize+chacha20poly1305.NonceSize)
	if keyGenErr != nil {
		return nil, nil, keyGenErr
	}

	chachaKey := secret[:chacha20poly1305.KeySize]
	nonce := secret[chacha20poly1305.KeySize:]

	return chachaKey, nonce, nil
}
