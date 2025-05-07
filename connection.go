package emailer

import (
	"io"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type Client interface {
	List(string, string, chan *imap.MailboxInfo) error
	Fetch(*imap.SeqSet, []imap.FetchItem, chan *imap.Message) error
	Select(string, bool) (*imap.MailboxStatus, error)
	Logout() error
}

type SMTPClient interface {
	SendMail(string, []string, io.Reader) error
}

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
