package emailer

import (
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
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
