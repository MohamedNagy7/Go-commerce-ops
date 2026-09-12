package email

import (
	"context"
	"os"

	"github.com/resend/resend-go/v4"
)

type ResendProvider struct {
	client *resend.Client
	from   string
}

func NewResendProvider() *ResendProvider {
	return &ResendProvider{
		client: resend.NewClient(os.Getenv("RESEND_API_KEY")),
		from:   os.Getenv("RESEND_FROM_EMAIL"),
	}
}

func (r *ResendProvider) Send(ctx context.Context, msg Message, fileContent string) error {
	_, err := r.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:    r.from,
		To:      []string{msg.To},
		Subject: msg.Subject,
		Html:    msg.Html,
		ReplyTo: r.from,

		Attachments: []*resend.Attachment{
			{
				Filename: "sample.txt",
				Content:  []byte(fileContent),
			},
		},
	})
	return err
}
