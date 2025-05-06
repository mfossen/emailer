package main

import (
	"context"
	"fmt"
	"io"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/mfossen/emailer"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"
)

func listMessages(ctx context.Context, cmd *cli.Command) error {
	client, err := newClient(cmd)
	if err != nil {
		return err
	}
	defer client.Logout()

	msgs, err := emailer.ListMessages(client, cmd.String("mailbox"))
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Date", "From", "Subject"})

	for _, m := range msgs {
		var from string
		for _, addr := range m.Envelope.From {
			from += addr.Address() + ", "
		}
		if len(m.Envelope.From) == 1 {
			from = strings.TrimSuffix(from, ", ")
		}
		table.Append([]string{strconv.FormatUint(uint64(m.SeqNum), 10), m.Envelope.Date.String(), from, m.Envelope.Subject})
	}

	table.Render()
	return nil
}

func showMessage(ctx context.Context, cmd *cli.Command) error {
	client, err := newClient(cmd)
	if err != nil {
		return err
	}
	defer client.Logout()

	msgID := cmd.Uint32Slice("id")

	msgs, err := emailer.GetMessages(client, cmd.String("mailbox"), msgID...)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		r := msg.GetBody(&imap.BodySectionName{})
		m, err := mail.ReadMessage(r)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(m.Body)
		if err != nil {
			return err
		}

		fmt.Println(string(body))

		var from string
		for _, addr := range msg.Envelope.From {
			from += addr.Address() + ", "
		}
		if len(msg.Envelope.From) == 1 {
			from = strings.TrimSuffix(from, ", ")
		}

		fmt.Printf("%T\n", msg)
		fmt.Printf("ID: %v\n", msg.SeqNum)
		fmt.Printf("Date: %s\n", msg.Envelope.Date)
		fmt.Printf("From: %v\n", from)
		fmt.Printf("Subject: %s\n", msg.Envelope.Subject)
		fmt.Printf("Body: %v\n", string(body))
	}

	return nil
}
