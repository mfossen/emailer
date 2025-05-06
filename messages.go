package emailer

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type fetch struct {
	from  uint32
	to    uint32
	items []uint32
}

func fetchItems(client *client.Client, fetch fetch) ([]*imap.Message, error) {
	msgs := []*imap.Message{}
	seqSet := new(imap.SeqSet)

	if fetch.from > 0 && fetch.to > 0 {
		seqSet.AddRange(fetch.from, fetch.to)
	} else {
		seqSet.AddNum(fetch.items...)
	}

	section := &imap.BodySectionName{}
	fetchItem := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- client.Fetch(seqSet, fetchItem, messages)
	}()

	for m := range messages {
		msgs = append(msgs, m)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return msgs, nil
}

func ListMessages(client *client.Client, mailbox string) ([]*imap.Message, error) {

	mbox, err := client.Select(mailbox, true)
	if err != nil {
		return nil, err
	}

	if mbox.Messages == 0 {
		return []*imap.Message{}, nil
	}

	from := uint32(1)
	to := mbox.Messages

	msgs, err := fetchItems(client, fetch{from: from, to: to})
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func GetMessages(client *client.Client, mailbox string, ids ...uint32) ([]*imap.Message, error) {

	_, err := client.Select(mailbox, true)
	if err != nil {
		return nil, err
	}

	msgs, err := fetchItems(client, fetch{items: ids})
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
