package emailer

import (
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// NewTLSClient creates a new authenticated TLS client
func NewTLSClient(auth sasl.Client, address string) (*client.Client, error) {
	client, err := client.DialTLS(address, nil)
	if err != nil {
		return nil, err
	}
	if err := client.Authenticate(auth); err != nil {
		return nil, err
	}
	return client, nil
}

func NewSMTPClient(auth sasl.Client, address string) (*smtp.Client, error) {
	client, err := smtp.DialTLS(address, nil)
	if err != nil {
		return nil, err
	}

	if err := client.Auth(auth); err != nil {
		return nil, err
	}

	return client, nil
}
