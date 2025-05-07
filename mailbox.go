package emailer

import (
	"github.com/emersion/go-imap"
)

// ListMailboxes collects all available mailboxes of an IMAP account
// and returns them in a slice
// TODO: implement limits, pagination, or splitting by folder level
// to handle cases of large or complex folder structures
func ListMailboxes(client Client) ([]*imap.MailboxInfo, error) {
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
