package emailer

import (
	"io"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// Client provides the necessary functions for listing and fetching
// IMAP mailboxes and messages.
// A go-imap *client.Client is a valid implementation.
type Client interface {
	List(string, string, chan *imap.MailboxInfo) error
	Fetch(*imap.SeqSet, []imap.FetchItem, chan *imap.Message) error
	Select(string, bool) (*imap.MailboxStatus, error)
	Logout() error
}

// SMTPClient provides the necessary functions for sending an email.
// A go-smtp *smtp.Client is a valid implementation.
type SMTPClient interface {
	SendMail(string, []string, io.Reader) error
	Quit() error
}

// NewTLSClient creates a new authenticated TLS client using Plain auth
// and returns a ready to use client for IMAP operations.
// The caller is responsible for calling client.Logout()
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

// NewSMTPClient creates a new authenticated TLS client using Plain auth
// and returns a ready to use client for SMTP operations.
// The caller is responsible for calling client.Quit()
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
