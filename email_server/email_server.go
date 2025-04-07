package email_server

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"sync"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"go.beyondstorage.io/v5/types"
)

type ImapFetcher struct {
	ImapClient *client.Client
	Email      string
	Password   string
	types.UnimplementedStorager
	mu      sync.Mutex
	child   map[string]*ImapFetcher
	workDir string
}

type EmailBuffer struct {
	FileName   string
	Attachment bytes.Buffer
}

type EmailAttachments struct {
	Attachments []EmailBuffer
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		result[i] = charset[randomIndex.Int64()]
	}
	return string(result)
}

func (e *EmailAttachments) Write(p []byte) (n int, err error) {
	var buf bytes.Buffer
	buf.Write(p)
	e.Attachments = append(e.Attachments, EmailBuffer{
		FileName:   generateRandomString(10),
		Attachment: buf,
	})

	// e.Attachments[0].Attachment.Write(p)
	// e.Attachments[0].FileName = "testfilerandom"

	return 0, nil
}

func (e *EmailAttachments) WriteWithName(filename string, p []byte) (n int, err error) {
	var buf bytes.Buffer
	buf.Write(p)
	e.Attachments = append(e.Attachments, EmailBuffer{
		FileName:   filename,
		Attachment: buf,
	})
	return 0, nil
}

func GetImapFetcher() (*ImapFetcher, error) {

	fetcher := &ImapFetcher{
		Email:    "randomfavone@gmail.com",
		Password: "hylc loju xhhb yqar",
	}

	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return nil, err
	}

	if err = c.Login("randomfavone@gmail.com", "hylc loju xhhb yqar"); err != nil {
		return nil, err
	}
	fetcher.ImapClient = c
	return fetcher, nil

}

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

	fmt.Println("Messages in the mailbox:", mBox.Messages)
	for i := uint32(mBox.Messages - 10); i <= mBox.Messages; i++ {

		fmt.Println("$$$$Fetching message:", i)

		seqSet := new(imap.SeqSet)
		seqSet.AddNum(i)

		section := &imap.BodySectionName{}
		items := []imap.FetchItem{section.FetchItem()}

		messages := make(chan *imap.Message, 1)
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

		fmt.Println("### email type : ", mr.Header.Get("Content-Type"))

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
	}

	fmt.Println("Done")

}

func (m *ImapFetcher) String() string {
	return fmt.Sprintf("Email: %s", m.Email)
}

func (m *ImapFetcher) Create(path string, pairs ...types.Pair) *types.Object {
	return &types.Object{}
}

func (m *ImapFetcher) Delete(path string, pairs ...types.Pair) error {
	return nil
}

func (m *ImapFetcher) DeleteWithContext(ctx context.Context, path string, pairs ...types.Pair) error {
	return nil
}

type test struct {
	message string
}

func (m test) ContinuationToken() string {
	return m.message
}

func (m *ImapFetcher) List(path string, pairs ...types.Pair) (*types.ObjectIterator, error) {
	ctx := context.Background()
	return types.NewObjectIterator(ctx, nil, nil), nil
}

func (m *ImapFetcher) ListWithContext(ctx context.Context, path string, pairs ...types.Pair) (*types.ObjectIterator, error) {
	return nil, nil
}

func (m *ImapFetcher) Metadata(pairs ...types.Pair) *types.StorageMeta {
	return nil
}

func (m *ImapFetcher) Read(path string, w io.Writer, pairs ...types.Pair) (int64, error) {
	mBox, err := m.ImapClient.Select("INBOX", false)
	if err != nil {
		panic(err)
	}
	if mBox.Messages == 0 {
		fmt.Println("No messages in the mailbox")
		return 0, fmt.Errorf("No messages in the mailbox")
	}

	fmt.Println("Messages in the mailbox:", mBox.Messages)
	for i := uint32(mBox.Messages - 50); i <= mBox.Messages; i++ {

		fmt.Println("Fetching email:", i)

		seqSet := new(imap.SeqSet)
		seqSet.AddNum(i)

		section := &imap.BodySectionName{}
		items := []imap.FetchItem{section.FetchItem()}

		messages := make(chan *imap.Message, 1)
		// go func() {
		if err := m.ImapClient.Fetch(seqSet, items, messages); err != nil {
			fmt.Println("Failed to fetch:", err)
		}
		// }()

		msg := <-messages
		if msg == nil {
			fmt.Println("No message found")
			return 0, fmt.Errorf("No message found")
		}

		r := msg.GetBody(section)
		if r == nil {
			fmt.Println("Failed to get body")
			return 0, fmt.Errorf("Failed to get body")
		}

		// Parse the message
		mr, err := mail.CreateReader(r)
		if err != nil {
			fmt.Println("Failed to create mail reader:", err)
			return 0, fmt.Errorf("Failed to create mail reader")
		}

		for {
			mail_body, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Failed to read part:", err)
				return 0, fmt.Errorf("Failed to read part")
			}

			switch h := mail_body.Header.(type) {
			case *mail.AttachmentHeader:
				filename, _ := h.Filename()
				os.Mkdir("attachments", os.ModePerm)

				fmt.Println("Found attachment: ", filename)

				// f, err := os.Create(fmt.Sprintf("attachments/%s", filename))
				// if err != nil {
				// 	return 0, fmt.Errorf("Failed to create file")
				// }
				// defer f.Close()

				// var buf bytes.Buffer
				data := make([]byte, 0, 1024)
				var totalBytes int = 0
				for {
					by := make([]byte, 1024)
					n, err := mail_body.Body.Read(by)
					fmt.Println("#### n : ", n)
					if err != nil {
						if err == io.EOF {
							break
						}
						return 0, err
					}
					data = append(data, by[:n]...)
					totalBytes += n
				}

				fmt.Println(filename, "totalBytes: ", totalBytes)
				fmt.Println(filename, "data : ", len(data))

				w.Write(data)

				fmt.Println("Attachment saved:", filename, "\n")
			}
		}
	}

	fmt.Println("Attachments Fetched Successfully")
	return 0, nil
}

func (m *ImapFetcher) ReadWithContext(ctx context.Context, path string, w io.Writer, pairs ...types.Pair) (int64, error) {
	return 0, nil
}

func (m *ImapFetcher) Stat(path string, pairs ...types.Pair) (*types.Object, error) {
	return &types.Object{}, nil
}

func (m *ImapFetcher) StatWithContext(ctx context.Context, path string, pairs ...types.Pair) (*types.Object, error) {
	return &types.Object{}, nil
}

func (m *ImapFetcher) Write(path string, r io.Reader, size int64, pairs ...types.Pair) (int64, error) {
	return 0, nil
}

func (m *ImapFetcher) WriteWithContext(ctx context.Context, path string, r io.Reader, size int64, pairs ...types.Pair) (int64, error) {
	return 0, nil
}

func (m *ImapFetcher) mustEmbedUnimplementedStorager() {}

func EmailFn(data ...types.Pair) (types.Storager, error) {
	fetcher, err := GetImapFetcher()
	if err != nil {
		panic(err)
	}
	return fetcher, nil
}
