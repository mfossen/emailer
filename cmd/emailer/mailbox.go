package main

import (
	"context"
	"os"
	"strings"

	"github.com/mfossen/emailer"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"
)

func listMailboxes(c context.Context, cmd *cli.Command) error {
	client, err := emailer.NewTLSClient(emailer.NewAuth(cmd.String("username"), cmd.String("password")))
	if err != nil {
		return err
	}
	defer client.Logout()

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
