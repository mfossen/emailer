package emailer

import (
	"github.com/emersion/go-sasl"
)

func NewAuth(username string, password string) sasl.Client {
	return sasl.NewPlainClient("", username, password)
}
