package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"secure-chat-client/client"
	"secure-chat-client/state"
	"secure-chat-client/util"
)

func Login(username, password string) error {
	c, cErr := util.CreateClient(nil)
	if cErr != nil {
		state.ClearAuthState()
		return errors.New("failed to create client")
	}

	body := client.LoginRequestBody{
		Username: username,
		Password: password,
	}

	res, resErr := c.PostAuthLoginWithResponse(context.Background(), body)
	if resErr != nil {
		state.ClearAuthState()
		return errors.New("failed to login")
	}

	if !util.Is2xx(res.HTTPResponse) {
		state.ClearAuthState()
		return errors.New("failed to login")
	}

	state.AccessToken = &res.JSON200.AccessToken
	state.RefreshToken = &res.JSON200.RefreshToken

	decodedPubKey, pubKeyErr := base64.StdEncoding.DecodeString(res.JSON200.PubKey)
	if pubKeyErr != nil {
		state.ClearAuthState()
		return pubKeyErr
	}
	state.PubKey = decodedPubKey

	encPrivKey, privKeyErr := base64.StdEncoding.DecodeString(res.JSON200.EncPrivKey)
	if privKeyErr != nil {
		state.ClearAuthState()
		return privKeyErr
	}
	decPrivKey, decPrivKeyErr := decryptPrivateKey(password, encPrivKey)
	if decPrivKeyErr != nil {
		state.ClearAuthState()
		return decPrivKeyErr
	}
	state.PrivKey = decPrivKey

	dataErr := state.FetchDataAfterLogin()
	if dataErr != nil {
		state.ClearAuthState()
		return dataErr
	}

	wsErr := state.ConnectWebSocket(*state.AccessToken)
	if wsErr != nil {
		state.ClearAuthState()
		return wsErr
	}
	return nil
}
