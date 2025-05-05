package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			&cli.Command{
				Name:    "mailbox",
				Aliases: []string{"mbox"},
				Usage:   "Operations on mailboxes",
				Commands: []*cli.Command{
					&cli.Command{
						Name:  "list",
						Usage: "List mailboxes",
						Action: func(c context.Context, cmd *cli.Command) error {
							fmt.Println("Listing mailboxes")
							return nil
						},
					},
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
