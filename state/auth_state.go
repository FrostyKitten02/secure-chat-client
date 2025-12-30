package state

var AccessToken *string = nil
var RefreshToken *string = nil

var PubKey []byte = nil
var PrivKey []byte = nil

func ClearAuthState() {
	AccessToken = nil
	RefreshToken = nil
	PubKey = nil
	PrivKey = nil
	ClearDataState()
}

func IsLoggedIn() bool {
	return AccessToken != nil
}
