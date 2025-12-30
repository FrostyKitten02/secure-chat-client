package service

import (
	"context"
	"errors"
	"secure-chat-client/client"
	"secure-chat-client/state"
	"secure-chat-client/util"
)

func AddUser(user client.UserListItemDto) error {
	c, cErr := util.CreateClient(state.AccessToken)
	if cErr != nil {
		return cErr
	}

	if user.Id == "" {
		return errors.New("error adding user, try again")
	}

	res, resErr := c.PostUsersbyUserIdAddToChat(context.Background(), user.Id)

	if resErr != nil {
		return resErr
	}

	if !util.Is2xx(res) {
		return errors.New("error adding user, try again")
	}

	return nil
}
