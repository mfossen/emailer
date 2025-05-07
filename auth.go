package emailer

import (
	"github.com/emersion/go-sasl"
)

// NewAuth returns a sasl Plain auth client to use for authenticating
// to IMAP and SMTP servers.
func NewAuth(username string, password string) sasl.Client {
	return sasl.NewPlainClient("", username, password)
}
