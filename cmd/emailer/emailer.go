package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

// main sets up command line subcommands and necessary flags before running the
// application at the end of the function.
// All of the actual logic is in adjacent files that contain ActionFunc's which are used as the
// cli.Command Action that will get called.
func main() {

	mailboxFlag := &cli.StringFlag{
		Name:     "mailbox",
		Aliases:  []string{"m"},
		Usage:    "Mailbox to operate on",
		Required: true,
	}

	cmd := &cli.Command{
		Usage: "command-line application for IMAP and SMTP operations",
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
				Commands: []*cli.Command{
					&cli.Command{
						Name:   "list",
						Usage:  "list emails in a mailbox",
						Action: listMessages,
						Flags:  []cli.Flag{mailboxFlag},
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
							mailboxFlag,
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
							&cli.StringFlag{
								Name:     "smtp-password",
								Usage:    "SMTP password to use",
								Required: true,
								Sources:  cli.EnvVars("SMTP_PASSWORD", "PASSWORD"),
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
