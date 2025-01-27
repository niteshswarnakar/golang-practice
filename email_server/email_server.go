package email_server

import (
	"fmt"
	"io"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

func EmailTest() {

	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		panic(err)
	}

	if err = c.Login("randomfavone@gmail.com", "hylc loju xhhb yqar"); err != nil {
		panic(err)
	}
	fmt.Println("Logged in")
	mBox, err := c.Select("INBOX", false)
	if err != nil {
		panic(err)
	}
	if mBox.Messages == 0 {
		fmt.Println("No messages in the mailbox")
		return
	}

	seqSet := new(imap.SeqSet)
	fmt.Println("Messages in the mailbox:", mBox.Messages)
	seqSet.AddNum(mBox.Messages)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, mBox.Messages)

	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			fmt.Println("Failed to fetch:", err)
		}
	}()

	msg := <-messages
	if msg == nil {
		fmt.Println("No message found")
		return
	}

	r := msg.GetBody(section)
	if r == nil {
		fmt.Println("Failed to get body")
		return
	}

	// Parse the message
	mr, err := mail.CreateReader(r)
	if err != nil {
		fmt.Println("Failed to create mail reader:", err)
		return
	}

	for {
		mail_body, err := mr.NextPart()
		if err == io.EOF {
			fmt.Println("====End of message")
			break
		}
		if err != nil {
			fmt.Println("Failed to read part:", err)
			return
		}

		switch h := mail_body.Header.(type) {
		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			fmt.Println("Found attachment:", filename)

			err := os.Mkdir("attachments", os.ModePerm)
			if err != nil {
				fmt.Println("Failed to create directory:", err)
			}

			fmt.Println("Saving attachment:", filename)

			f, err := os.Create(fmt.Sprintf("attachments/%s", filename))
			if err != nil {
				fmt.Println("Failed to create file:", err)
				return
			}
			defer f.Close()
			if _, err := io.Copy(f, mail_body.Body); err != nil {
				fmt.Println("Failed to save attachment:", err)
				return
			}
			fmt.Println("Attachment saved:", filename)
		}
	}

	fmt.Println("Done")

}
