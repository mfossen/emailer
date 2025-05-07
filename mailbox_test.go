package emailer

import (
	"testing"

	"github.com/emersion/go-imap"
	"github.com/stretchr/testify/assert"
)

type MockClient struct{}

func (m MockClient) List(ref string, name string, ch chan *imap.MailboxInfo) error {
	mailboxInfo := []*imap.MailboxInfo{
		&imap.MailboxInfo{
			Name: "INBOX",
		},
		&imap.MailboxInfo{
			Name: "Junk",
		},
		&imap.MailboxInfo{
			Name: "Trash",
		},
	}
	for _, info := range mailboxInfo {
		ch <- info
	}
	close(ch)
	return nil
}

func (m MockClient) Select(name string, readOnly bool) (*imap.MailboxStatus, error) {
	mbox := &imap.MailboxStatus{
		Name:     name,
		ReadOnly: readOnly,
		Messages: 3,
	}
	return mbox, nil
}

func (m MockClient) Logout() error {
	return nil
}

func (m MockClient) Fetch(seqnums *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	defer close(ch)
	msgs := []*imap.Message{
		&imap.Message{
			SeqNum: 1,
		},
		&imap.Message{
			SeqNum: 2,
		},
		&imap.Message{
			SeqNum: 3,
		},
	}
	for _, msg := range msgs {
		if seqnums.Contains(msg.SeqNum) {
			ch <- msg
		}
	}
	return nil
}

func TestListMailbox(t *testing.T) {
	// Test cases for the Mailbox struct
	mockClient := new(MockClient)

	mailboxInfo, err := ListMailboxes(mockClient)

	assert.NotEmpty(t, mailboxInfo)
	assert.Len(t, mailboxInfo, 3)
	assert.NoError(t, err)
}
