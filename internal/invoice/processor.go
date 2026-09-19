package invoice

import (
	"context"
	"fmt"

	"github.com/MohamedNagy7/Go-commerce-ops/internal/email"
	"github.com/MohamedNagy7/Go-commerce-ops/internal/rabbitmq"
)

type PDFGenerator func(InvoiceData) ([]byte, error)

type Processor struct {
	GeneratePDF PDFGenerator
	Email       email.Provider
}

func (p *Processor) Process(ctx context.Context, payload rabbitmq.InvoiceRequestedPayload) error {
	data := BuildInvoiceData(payload)

	pdfBytes, err := p.GeneratePDF(data)

	if err != nil {
		return err
	}

	if err := p.Email.Send(ctx, email.Message{
		To:      data.CustomerEmail,
		Subject: "Your Order From Drip",
		Html:    fmt.Sprintf("<p>Invoice for order %s attached.</p>", data.OrderID),
	}, pdfBytes, data.OrderID+".pdf"); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}
