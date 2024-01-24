package wrapper

import (
	"github.com/rs/zerolog"
	"github.com/zalando/go-keyring"
)

type KeyRing struct {
	Logger *zerolog.Logger
}

func (k *KeyRing) Set(app, user, pass string) bool {
	err := keyring.Set(app, user, pass)
	if err != nil {
		k.Logger.Error().AnErr("KeyRing::Set::Error", err)
		return false
	}

	return true
}

func (k *KeyRing) Get(app, user string) (string, bool) {
	passwd, err := keyring.Get(app, user)
	if err != nil {
		k.Logger.Error().AnErr("KeyRing::Get::Error", err)
		return "", false
	}
	return passwd, true
}
