package email

import "context"

type Message struct {
	To      string
	Subject string
	Html    string
}

type Provider interface {
	Send(ctx context.Context, msg Message, attachment []byte, filename string) error
}
