package main

import (
	"github.com/emersion/go-imap/client"
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
