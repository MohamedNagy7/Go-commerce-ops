package invoice

import (
	"context"
	"fmt"

	"github.com/MohamedNagy7/Go-commerce-ops/internal/email"
	"github.com/MohamedNagy7/Go-commerce-ops/internal/rabbitmq"
)

type PDFGenerator func(InvoiceData, []byte) ([]byte, error)

type Processor struct {
	GeneratePDF PDFGenerator
	Email       email.Provider
	LogoPNG     []byte
}

func (p *Processor) Process(ctx context.Context, payload rabbitmq.InvoiceRequestedPayload) error {
	data := BuildInvoiceData(payload)

	pdfBytes, err := p.GeneratePDF(data, p.LogoPNG)
	if err != nil {
		return fmt.Errorf("generate pdf: %w", err)
	}

	if err := p.Email.Send(ctx, email.Message{
		To:      data.CustomerEmail,
		Subject: fmt.Sprintf("Your invoice — Order #%s", data.OrderID),
		Html:    BuildInvoiceHTML(data),
		Attachments: []email.Attachment{
			{Content: pdfBytes, Filename: data.OrderID + ".pdf", ContentType: "application/pdf"},
			{Content: p.LogoPNG, Filename: "logo.png", ContentType: "image/png", ContentID: "logo"},
		},
	}); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}
