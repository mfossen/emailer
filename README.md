# emailer: cli IMAP & SMTP

[![builds.sr.ht status](https://builds.sr.ht/~mfossen/emailer.svg)](https://builds.sr.ht/~mfossen/emailer?)

## Usage

`./emailer --help` for a list of commands and flags available, and `./emailer <command> --help` for further help info.

### Commands

For the following commands, these flags will need to be set:

- `--username` or `USERNAME` environment variable: IMAP username
- `--password` or `PASSWORD` environment variable: IMAP password
- `--address` or `IMAP_ADDRESS` environment variable: IMAP server address, in the format of `imap.example.com:993` (note: only TLS is supported)

---

`emailer mailbox list`: list available mailboxes

```
+---------+------------+
|  NAME   | ATTRIBUTES |
+---------+------------+
| Archive | \Archive   |
| Drafts  | \Drafts    |
| INBOX   |            |
| Junk    | \Junk      |
| Sent    | \Sent      |
| Trash   | \Trash     |
+---------+------------+
```

---

The next commands operate on messages in a mailbox, so require `--mailbox <mailbox name>` to be given.

- `emailer message --mailbox INBOX list`

```
+----+--------------------------------+------------------------+------------------------+
| ID |              DATE              |          FROM          |        SUBJECT         |
+----+--------------------------------+------------------------+------------------------+
|  1 | 2025-05-06 17:23:37 +0000      | support@example.com    | Welcome!               |
|    | +0000                          |                        |                        |
|  2 | 2025-05-06 12:23:50 -0500 CDT  | tester@gmail.com       | test email             |
|  3 | 2025-05-06 17:54:32 -0500 CDT  | foo@example.com        | test another email     |
+----+--------------------------------+------------------------+------------------------+
```


- `emailer message --mailbox INBOX show --id 3 --id 2` will attempt to open the messages given by the `id` flags in `$PAGER`, falling back to `less`, and finally just printing out if neither of those are available.

```
ID: 2
Date: 2025-05-06 12:23:50 -0500 CDT
From: tester@gmail.com
Subject: test email

test email message


ID: 3
Date: 2025-05-06 17:54:32 -0500 CDT
From: foo@example.com
Subject: test another email

test another email message

one
two
three
```

---

The final supported command is for sending an email, with the required flags:

- `--smtp-address` or `SMTP_ADDRESS` environment variable: address of SMTP server to use in the format of `smtp.example.com:465` (note: only TLS is supported)
- `--smtp-username` or one of `SMTP_USERNAME` or `USERNAME`: SMTP username
- `--smtp-password` or one of `SMTP_PASSWORD` or `PASSWORD`: SMTP password

- `emailer message send` will attempt to open `$EDITOR`, falling back to `vim` with the content of a simple message filled out. On saving and closing, the message will be sent.

```
From: <test@example.com>
To:
Subject:

<!--- enter body text below this line (this will get removed before sending) --->

```

## Contributing

### Quickstart: building and testing

- building: `go build -mod=vendor ./cmd/...`
- testing: `go test -mod=vendor ./...`

CI builds are available on [builds.sr.ht/~mfossen/emailer](https://builds.sr.ht/~mfossen/emailer) which run tests, validates the binary will run, and runs builds for Linux, Windows, and Darwin `GOOS` targets as a sanity check.


### Dependencies

The primary outside dependencies this project relies on are:

- [urfave/cli](https://github.com/urfave/cli) for setting up the CLI flag and actions bits
- [emersion/go-imap](https://github.com/emersion/go-imap/tree/v1) (note: v1 version is used) for the IMAP bits
- [emersion/go-smtp](https://github.com/emersion/go-smtp) for the SMTP bits
- [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) for outputting nicer looking tables


### Project Structure

Top-level `*.go` files contain exported functions and corresponding logic for running operations and collecting returned data into nicer types to use by the caller e.g. selecting a mailbox, listing messages, and collecting them from a channel into a slice that gets returned.

Files in `./cmd/emailer` are all part of the `main` package, with `emailer.go` wiring up the command-line bits and separate files setting up compatible [urfave/cli ActionFunc](https://pkg.go.dev/github.com/urfave/cli/v3#ActionFunc) that get called.

Mostly, these functions set up a client, call the exported functions from the top-level, and handle the input/output of the returned content.

### Improvements and Next Steps

- implementing remaining functionality people would reasonably expect, available in [go-imap](https://pkg.go.dev/github.com/emersion/go-imap/client) such as searching, renames, delete, etc.
- nicer user experience when sending emails, highest priority would be responding to an existing message from the IMAP side, followed by things like guided email creation, handling multipart emails, and dealing with attachments.
- addressing potential performance issues: listing mailboxes and messages aren't limited or paginated in any way, so it's likely that under large numbers of emails (or complex folder setups) problems will crop up with fetching and outputting those.
