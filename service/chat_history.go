package service

import (
	"context"
	"log/slog"
	"secure-chat-client/model"
	"secure-chat-client/state"
	"secure-chat-client/util"
)

func GetCurrentChatHistory() ([]model.DecryptedDirectMessage, error) {
	c, cErr := util.CreateClient(state.AccessToken)
	if cErr != nil {
		return nil, cErr
	}

	res, resErr := c.GetChatsHistoryByUserByUserIdWithResponse(context.Background(), state.CurrentChatUserId)
	if resErr != nil {
		return nil, resErr
	}

	if res.JSON200.DirectMessages == nil {
		return nil, nil
	}

	if res.JSON200.FromUser == nil || len(res.JSON200.FromUser) == 0 {
		return nil, nil
	}

	fromUser := res.JSON200.FromUser[0]
	msgs := make([]model.DecryptedDirectMessage, len(res.JSON200.DirectMessages))
	for i, dm := range res.JSON200.DirectMessages {
		var username string = ""
		if dm.FromUserId == fromUser.UserId {
			username = fromUser.Username
		} else {
			username = "You"
		}

		plainText, decErr := state.DecryptMessage(dm.Nonce, dm.CipherText)
		if decErr != nil {
			slog.Error("error decrypting message", dm.Nonce, dm.CipherText, decErr)
			msgs[i] = model.DecryptedDirectMessage{
				FromUsername: username,
				PlainText:    dm.CipherText,
			}
			continue
		}

		msgs[i] = model.DecryptedDirectMessage{
			FromUsername: username,
			PlainText:    plainText,
		}
	}

	return msgs, nil
}
