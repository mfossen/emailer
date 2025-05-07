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
						Name:   "list",
						Usage:  "List mailboxes",
						Action: listMailboxes,
					},
				},
			},
			&cli.Command{
				Name:    "message",
				Aliases: []string{"msg"},
				Usage:   "Operations on email messages",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "mailbox",
						Aliases:  []string{"m"},
						Usage:    "Mailbox to operate on",
						Required: false,
					},
				},
				Commands: []*cli.Command{
					&cli.Command{
						Name:   "list",
						Usage:  "list emails in a mailbox",
						Action: listMessages,
					},
					&cli.Command{
						Name:   "show",
						Usage:  "show email messages",
						Action: showMessage,
						Flags: []cli.Flag{
							&cli.Uint32SliceFlag{
								Name:     "id",
								Usage:    "message id to show",
								Required: true,
							},
						},
					},
					&cli.Command{
						Name:   "send",
						Usage:  "send an email",
						Action: sendMessage,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "smtp-address",
								Usage:    "TLS SMTP server address to use",
								Required: true,
								Sources:  cli.EnvVars("SMTP_ADDRESS"),
							},
							&cli.StringFlag{
								Name:     "smtp-username",
								Usage:    "SMTP username to use",
								Required: true,
								Sources:  cli.EnvVars("SMTP_USERNAME", "USERNAME"),
							},
						},
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Usage:    "IMAP username",
				Required: true,
				Sources:  cli.EnvVars("USERNAME"),
			},
			&cli.StringFlag{
				Name:     "password",
				Usage:    "IMAP password",
				Required: true,
				Sources:  cli.EnvVars("PASSWORD"),
			},
			&cli.StringFlag{
				Name:     "address",
				Usage:    "TLS IMAP server address",
				Required: true,
				Sources:  cli.EnvVars("IMAP_ADDRESS"),
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
