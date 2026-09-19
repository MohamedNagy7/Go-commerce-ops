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

func (r *ResendProvider) Send(ctx context.Context, msg Message) error {
	atts := make([]*resend.Attachment, 0, len(msg.Attachments))
	for _, a := range msg.Attachments {
		atts = append(atts, &resend.Attachment{
			Content:     a.Content,
			Filename:    a.Filename,
			ContentType: a.ContentType,
			ContentId:   a.ContentID,
		})
	}
	_, err := r.client.Emails.SendWithContext(ctx, &resend.SendEmailRequest{
		From:        r.from,
		To:          []string{msg.To},
		Subject:     msg.Subject,
		Html:        msg.Html,
		ReplyTo:     r.from,
		Attachments: atts,
	})
	return err
}
