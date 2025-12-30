package service

import (
	"context"
	"secure-chat-client/client"
	"secure-chat-client/state"
	"secure-chat-client/util"
)

func Search(query string) ([]client.UserListItemDto, error) {
	c, cErr := util.CreateClient(state.AccessToken)
	if cErr != nil {
		return nil, cErr
	}

	res, resErr := c.GetUsersListWithResponse(context.Background(), &client.GetUsersListParams{Query: &query})
	if resErr != nil {
		return nil, resErr
	}

	if res.JSON200 == nil || res.JSON200.Users == nil {
		return []client.UserListItemDto{}, nil
	}

	return res.JSON200.Users, nil
}
