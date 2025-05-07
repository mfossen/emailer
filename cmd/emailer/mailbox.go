package main

import (
	"context"
	"os"
	"strings"

	"github.com/mfossen/emailer"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"
)

// listMailboxes calls the top-level ListMailboxes function to output a table
// of all available IMAP mailboxes with their Name and Attributes
func listMailboxes(c context.Context, cmd *cli.Command) error {
	client, err := newClient(cmd)
	if err != nil {
		return err
	}
	defer client.Logout() //nolint:errcheck

	mailboxes, err := emailer.ListMailboxes(client)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Attributes"})

	for _, mailbox := range mailboxes {
		table.Append([]string{mailbox.Name, strings.Join(mailbox.Attributes, ", ")})
	}

	table.Render()

	return nil
}
