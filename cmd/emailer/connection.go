package main

import (
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-smtp"
	"github.com/mfossen/emailer"
	"github.com/urfave/cli/v3"
)

func newClient(cmd *cli.Command) (*client.Client, error) {
	auth := emailer.NewAuth(cmd.String("username"), cmd.String("password"))
	client, err := emailer.NewTLSClient(auth, cmd.String("address"))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func newSMTPClient(cmd *cli.Command) (*smtp.Client, error) {
	auth := emailer.NewAuth(cmd.String("username"), cmd.String("password"))
	client, err := emailer.NewSMTPClient(auth, cmd.String("smtp-address"))
	if err != nil {
		return nil, err
	}
	return client, err
}
