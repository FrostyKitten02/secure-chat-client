package state

import (
	"context"
	"secure-chat-client/client"
	"secure-chat-client/util"
)

var Chats []client.ChatDto = []client.ChatDto{}

func FetchDataAfterLogin() error {
	c, err := util.CreateClient(AccessToken)
	if err != nil {
		return err
	}

	chatErr := fetchChats(c)
	if chatErr != nil {
		return chatErr
	}

	return nil
}

func fetchChats(c *client.ClientWithResponses) error {
	res, resErr := c.GetChatsWithResponse(context.Background())
	if resErr != nil {
		return resErr
	}

	resBody := res.JSON200
	if resBody == nil {
		return nil
	}

	Chats = resBody.Chats

	return nil
}

func ClearDataState() {
	Chats = []client.ChatDto{}
	CurrentChatUserId = ""
	currentChatKey = []byte{}
}
