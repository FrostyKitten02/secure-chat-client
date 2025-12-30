package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"secure-chat-client/state"
)

func SendMessage(userId, msg string) error {
	userIdUuid, uuidErr := uuid.Parse(userId)
	if uuidErr != nil {
		return uuidErr
	}

	cipherText, nonce, encErr := state.EncryptCurrentChat(userId, msg)
	if encErr != nil {
		return encErr
	}

	newMessage := state.WsSendNewMessage{
		ToUserID:           userIdUuid,
		CipherText:         cipherText,
		Nonce:              nonce,
		SenderIdentityId:   "NOT USED!",
		ReceiverIdentityId: "NOT USED!",
	}

	wsMsg, marshalErr := json.Marshal(newMessage)
	if marshalErr != nil {
		return marshalErr
	}

	state.SendMessage(wsMsg)

	return nil
}
