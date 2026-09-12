package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

const (
	Exchange = "app.events"
	DLX      = "app.events.dlx"

	InvoiceQueue      = "invoices.generate"
	InvoiceDLQ        = "invoices.generate.dlq"
	InvoiceRoutingKey = "invoice.requested"
)

// Idempotent — same guarantee as the Nest setupTopology. Safe to call on every boot
func SetupTopology(ch *amqp.Channel) error {
	if err := ch.ExchangeDeclare(Exchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.ExchangeDeclare(DLX, "topic", true, false, false, false, nil); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(InvoiceDLQ, true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind(InvoiceDLQ, InvoiceRoutingKey, DLX, false, nil); err != nil {
		return err
	}

	_, err := ch.QueueDeclare(InvoiceQueue, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    DLX,
		"x-dead-letter-routing-key": InvoiceRoutingKey,
	})
	if err != nil {
		return err
	}
	return ch.QueueBind(InvoiceQueue, InvoiceRoutingKey, Exchange, false, nil)
}
