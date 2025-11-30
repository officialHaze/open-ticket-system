package helper

import (
	"crypto/rand"
	"fmt"
	"ots/settings"

	"github.com/o1egl/paseto"
)

var Token *PasetoToken

func GeneratePaseto() error {
	b32key := make([]byte, 32)
	_, err := rand.Read(b32key)
	if err != nil {
		return fmt.Errorf("error generating 32-Byte random key")
	}

	Token = &PasetoToken{
		key:    b32key,
		footer: settings.MySettings.Get_TokenFooter(),
	}

	return nil
}

type PasetoToken struct {
	key    []byte
	footer string
}

func (p *PasetoToken) CreateToken(payload any) (string, error) {
	return paseto.NewV2().Encrypt(p.key, payload, p.footer)
}

func (p *PasetoToken) DecryptToken(token string, payload, footer any) error {
	return paseto.NewV2().Decrypt(token, p.key, payload, footer)
}
