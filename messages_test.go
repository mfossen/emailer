package emailer

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchItems(t *testing.T) {
	mockClient := MockClient{}

	msgs, err := fetchItems(mockClient, fetch{})
	assert.NoError(t, err)
	assert.Empty(t, msgs)

	msgs, err = fetchItems(mockClient, fetch{from: 1, to: 3})
	assert.NoError(t, err)
	assert.Len(t, msgs, 3)
}

func TestListMessages(t *testing.T) {
	mockClient := MockClient{}

	msgs, err := ListMessages(mockClient, "TEST")

	assert.NoError(t, err)
	assert.Len(t, msgs, 3)
}

func TestGetMessages(t *testing.T) {
	mockClient := MockClient{}

	msgs, err := GetMessages(mockClient, "TEST", 1)

	assert.NoError(t, err)
	assert.Len(t, msgs, 1)

	msgs, err = GetMessages(mockClient, "TEST", 1, 2)

	assert.NoError(t, err)
	assert.Len(t, msgs, 2)
}

type MockSMTPClient struct{}

func (c MockSMTPClient) SendMail(to string, from []string, r io.Reader) error {
	return nil
}

func (c MockSMTPClient) Quit() error {
	return nil
}

func TestSendMessage(t *testing.T) {
	mockClient := MockSMTPClient{}

	testMessages := []struct {
		from    string
		to      string
		wantErr bool
	}{
		{
			from:    "Good Address <good@example.com>",
			to:      "Good To <good@example.com>",
			wantErr: false,
		},
		{
			from:    "Good Address <good@example.com>",
			to:      "Good To One <good@example.com>, Good To Two <alsogood@example.com>",
			wantErr: false,
		},
		{
			from:    "Bad Address <bad>",
			to:      "Good To <good@example.com>",
			wantErr: true,
		},
		{
			from:    "Good Address <good@example.com>",
			to:      "Bad To <bad>",
			wantErr: true,
		},
		{
			from:    "bad",
			to:      "bad",
			wantErr: true,
		},
	}

	for _, msg := range testMessages {
		t.Run(msg.from+","+msg.to, func(tt *testing.T) {

			msgFormat := fmt.Sprintf("From: %s\nTo: %s\nSubject: test subject\n\n", msg.from, msg.to)
			err := SendMessage(mockClient, []byte(msgFormat))

			if msg.wantErr {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
		})
	}

	err := SendMessage(mockClient, []byte{})
	assert.Error(t, err)

}
