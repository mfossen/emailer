package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/mail"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/mfossen/emailer"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v3"
)

// listMessages gets all available messages from an IMAP mailbox and
// outputs them in a table with ID, Date, From, and Subject
func listMessages(ctx context.Context, cmd *cli.Command) error {
	client, err := newClient(cmd)
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

// showMessage fetches one or more messages from an IMAP mailbox and attempts to output them
// in the system pager given by the PAGER environment variable, "less" if PAGER isn't set, or
// finally just printing to standard out
func showMessage(ctx context.Context, cmd *cli.Command) error {
	client, err := newClient(cmd)
	if err != nil {
		return err
	}
	defer client.Logout()

	msgID := cmd.Uint32Slice("id")

	msgs, err := emailer.GetMessages(client, cmd.String("mailbox"), msgID...)
	if err != nil {
		return err
	}

	var pager string
	for _, attemptedPager := range []string{os.Getenv("PAGER"), "less"} {
		path, found := exec.LookPath(attemptedPager)
		if found == nil {
			pager = path
			break
		}
	}

	for _, msg := range msgs {
		r := msg.GetBody(&imap.BodySectionName{})
		m, err := mail.ReadMessage(r)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(m.Body)
		if err != nil {
			return err
		}

		var from string
		for _, addr := range msg.Envelope.From {
			from += addr.Address() + ", "
		}
		if len(msg.Envelope.From) == 1 {
			from = strings.TrimSuffix(from, ", ")
		}

		msgContent := strings.Builder{}

		_, err = msgContent.WriteString(fmt.Sprintf("ID: %d\n", msg.SeqNum))
		if err != nil {
			return err
		}
		_, err = msgContent.WriteString(fmt.Sprintf("Date: %s\n", msg.Envelope.Date))
		if err != nil {
			return err
		}
		_, err = msgContent.WriteString(fmt.Sprintf("From: %s\n", from))
		if err != nil {
			return err
		}
		_, err = msgContent.WriteString(fmt.Sprintf("Subject: %s\n", msg.Envelope.Subject))
		if err != nil {
			return err
		}
		_, err = msgContent.WriteString(fmt.Sprintf("\n%s\n", string(body)))
		if err != nil {
			return err
		}

		if pager == "" {
			fmt.Println(msgContent.String())
			continue
		}

		command := exec.Command(pager)
		command.Stdin = strings.NewReader(msgContent.String())
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			return err
		}
	}

	return nil
}

// sendMessage opens an available editor given by EDITOR or "vim" with a simple pre-filled
// set of email headers + a line marking where to enter the body text and sends the email on
// a save and close of the editor
func sendMessage(ctx context.Context, cmd *cli.Command) error {
	client, err := newSMTPClient(cmd)
	if err != nil {
		return err
	}
	defer client.Quit()

	msgTemplate := fmt.Sprintf(`From: <%s>
To:
Subject:

<!--- enter body text below this line (this will get removed before sending) --->
`, cmd.String("smtp-username"))

	tempfile, err := os.CreateTemp("", "emailer*")
	if err != nil {
		return err
	}

	err = os.WriteFile(tempfile.Name(), []byte(msgTemplate), 644)
	if err != nil {
		return err
	}
	defer os.Remove(tempfile.Name())

	var editor string
	for _, attemptedEditor := range []string{os.Getenv("EDITOR"), "vim"} {
		path, found := exec.LookPath(attemptedEditor)
		if found == nil {
			editor = path
			break
		}
	}

	if editor == "" {
		return fmt.Errorf("unable to open a file editor, set EDITOR environment variable")
	}

	command := exec.Command(editor, tempfile.Name())
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return err
	}

	msg, err := os.ReadFile(tempfile.Name())
	if err != nil {
		return err
	}

	reg, err := regexp.Compile(`<!---.*--->\n`)
	if err != nil {
		return err
	}
	foundBytes := reg.Find(msg)
	newMsg := bytes.Replace(msg, foundBytes, []byte(""), 1)
	err = emailer.SendMessage(client, newMsg)
	return err
}
