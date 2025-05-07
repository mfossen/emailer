package emailer

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"strings"

	"github.com/emersion/go-imap"
)

type fetch struct {
	from  uint32
	to    uint32
	items []uint32
}

type Message struct {
	From string
	To   string
	Msg  string
}

func fetchItems(client Client, fetch fetch) ([]*imap.Message, error) {
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

func ListMessages(client Client, mailbox string) ([]*imap.Message, error) {

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

func GetMessages(client Client, mailbox string, ids ...uint32) ([]*imap.Message, error) {

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

func SendMessage(client SMTPClient, msg []byte) error {
	msgReader := bytes.NewReader(msg)
	m, err := mail.ReadMessage(msgReader)
	if err != nil {
		return err
	}

	headers := m.Header

	addrParser := mail.AddressParser{}
	from, err := addrParser.Parse(headers.Get("From"))
	if err != nil {
		return err
	}

	to, err := addrParser.ParseList(headers.Get("To"))
	if err != nil {
		return err
	}

	if from == nil || len(to) == 0 {
		return fmt.Errorf("error parsing from or to addresses, got from:%v\nto:%v", from, to)
	}

	var toAddrs []string
	for _, a := range to {
		toAddrs = append(toAddrs, strings.Trim(a.Address, "<>"))
	}

	_, err = msgReader.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	return client.SendMail(strings.Trim(from.Address, "<>"), toAddrs, msgReader)
}
