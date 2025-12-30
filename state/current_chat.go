package state

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
	"io"
	"secure-chat-client/client"
)

var currentChatKey []byte = []byte{}
var CurrentChatUserId string = ""

// returns encrypted msg and nonce
func EncryptCurrentChat(userId, msg string) (string, string, error) {
	if userId != CurrentChatUserId {
		return "", "", errors.New("userId not match")
	}

	aead, err := chacha20poly1305.New(currentChatKey)
	if err != nil {
		return "", "", err
	}

	nonceSize := aead.NonceSize()
	nonce := make([]byte, nonceSize)
	_, nonceErr := rand.Read(nonce)
	if nonceErr != nil {
		return "", "", nonceErr
	}

	enc := aead.Seal(nil, nil, []byte(msg), nil)
	return base64.StdEncoding.EncodeToString(enc), base64.StdEncoding.EncodeToString(nonce), nil
}

func SetCurrentChat(chat client.ChatDto) error {
	//resetting current chat info
	CurrentChatUserId = ""
	currentChatKey = []byte{}

	if chat.User.PubKey == "" {
		return fmt.Errorf("no public key provided")
	}

	if chat.User.UserId == "" {
		return fmt.Errorf("no user id provided")
	}

	decodedPubKey, pubKeyErr := base64.StdEncoding.DecodeString(chat.User.PubKey)
	if pubKeyErr != nil {
		return pubKeyErr
	}

	key, keyErr := createChatKey(decodedPubKey, PrivKey)
	if keyErr != nil {
		return keyErr
	}

	CurrentChatUserId = chat.User.UserId
	currentChatKey = key
	return nil
}

func createChatKey(pubKey []byte, privKey []byte) ([]byte, error) {
	curve := ecdh.X25519()
	cPubKey, pubErr := curve.NewPublicKey(pubKey)
	if pubErr != nil {
		return nil, pubErr
	}

	cPrivKey, privErr := curve.NewPrivateKey(privKey)
	if privErr != nil {
		return nil, privErr
	}

	sharedSecret, sErr := cPrivKey.ECDH(cPubKey)
	if sErr != nil {
		return nil, sErr
	}

	derivedChatKey := hkdf.New(
		sha256.New,
		sharedSecret,
		nil,
		nil,
	)

	key := make([]byte, 32)
	_, err := io.ReadFull(derivedChatKey, key)
	return nil, err
}
