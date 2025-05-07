package emailer

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type Client interface {
	List(ref string, name string, ch chan *imap.MailboxInfo) error
	Fetch(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error
	Select(name string, readOnly bool) (*imap.MailboxStatus, error)
	Logout() error
}

type IMAPClient struct {
	client *client.Client
}

func (c IMAPClient) Logout() error {
	return c.client.Logout()
}

func (c IMAPClient) List(ref string, name string, ch chan *imap.MailboxInfo) error {
	return c.client.List(ref, name, ch)
}

func (c IMAPClient) Fetch(seqSet *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	return c.client.Fetch(seqSet, items, ch)
}

func (c IMAPClient) Select(name string, readOnly bool) (*imap.MailboxStatus, error) {
	return c.client.Select(name, readOnly)
}

// NewTLSClient creates a new authenticated TLS client
func NewTLSClient(auth sasl.Client, address string) (IMAPClient, error) {
	client, err := client.DialTLS(address, nil)
	if err != nil {
		return IMAPClient{}, err
	}
	if err := client.Authenticate(auth); err != nil {
		return IMAPClient{}, err
	}
	return IMAPClient{client}, nil
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
