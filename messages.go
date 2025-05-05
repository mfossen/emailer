package emailer

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func ListMessages(client *client.Client, mailbox string) ([]*imap.Message, error) {
	msgs := []*imap.Message{}

	mbox, err := client.Select(mailbox, true)
	if err != nil {
		return nil, err
	}

	from := uint32(1)
	to := mbox.Messages

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)
	items := []imap.FetchItem{imap.FetchEnvelope}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- client.Fetch(seqSet, items, messages)
	}()

	for m := range messages {
		msgs = append(msgs, m)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return msgs, nil
}
