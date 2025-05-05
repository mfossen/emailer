package emailer

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func ListMailboxes(client *client.Client) ([]*imap.MailboxInfo, error) {
	mboxes := []*imap.MailboxInfo{}

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		done <- client.List("", "*", mailboxes)
	}()

	for m := range mailboxes {
		mboxes = append(mboxes, m)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return mboxes, nil
}
