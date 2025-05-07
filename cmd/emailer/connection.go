package main

import (
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-smtp"
	"github.com/mfossen/emailer"
	"github.com/urfave/cli/v3"
)

// newClient is a helper function to return a new authenticated IMAP client
func newClient(cmd *cli.Command) (*client.Client, error) {
	auth := emailer.NewAuth(cmd.String("username"), cmd.String("password"))
	client, err := emailer.NewTLSClient(auth, cmd.String("address"))
	if err != nil {
		return client, err
	}
	return client, nil
}

// newSMTPClient is a helper function to return a new authenticated SMTP client
func newSMTPClient(cmd *cli.Command) (*smtp.Client, error) {
	auth := emailer.NewAuth(cmd.String("smtp-username"), cmd.String("smtp-password"))
	client, err := emailer.NewSMTPClient(auth, cmd.String("smtp-address"))
	if err != nil {
		return nil, err
	}
	return client, err
}
