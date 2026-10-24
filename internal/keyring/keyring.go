package keyring

import "github.com/zalando/go-keyring"

const (
	keyringService = "podlet"
	keyringUser    = "session-token"
)

func SaveToken(token string) error {
	return keyring.Set(keyringService, keyringUser, token)
}

func LoadToken() (string, error) {
	return keyring.Get(keyringService, keyringUser)
}

func DeleteToken() error {
	return keyring.Delete(keyringService, keyringUser)
}
