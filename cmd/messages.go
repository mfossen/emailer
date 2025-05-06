package main

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/mfossen/emailer"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"
)

func listMessages(ctx context.Context, cmd *cli.Command) error {
	client, err := emailer.NewTLSClient(emailer.NewAuth(cmd.String("username"), cmd.String("password")))
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
