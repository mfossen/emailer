package emailer

import (
	"os"

	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
)

// NewTLSClient creates a new authenticated TLS client
func NewTLSClient(auth sasl.Client) (*client.Client, error) {
	client, err := client.DialTLS(os.Getenv("IMAP_ADDRESS"), nil)
	if err != nil {
		return nil, err
	}
	if err := client.Authenticate(auth); err != nil {
		return nil, err
	}
	return client, nil
}
