package email

import "context"

type Attachment struct {
	Content     []byte
	Filename    string
	ContentType string
	ContentID   string // set to embed inline and reference via cid: in HTML
}

type Message struct {
	To          string
	Subject     string
	Html        string
	Attachments []Attachment
}

type Provider interface {
	Send(ctx context.Context, msg Message) error
}
