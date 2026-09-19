package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"runtime"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ProcessFunc is swapped out later for real PDF generation + upload.
type ProcessFunc func(ctx context.Context, payload InvoiceRequestedPayload) error

func StartInvoiceConsumer(ch *amqp.Channel, process ProcessFunc) error {
	// mirrors prefetch(1) from the Nest consumer, but here it caps how many
	// unacked messages RabbitMQ will hand us — the actual concurrency cap
	// is the semaphore below, sized to CPU cores since PDF rendering is CPU-bound
	if err := ch.Qos(runtime.NumCPU(), 0, false); err != nil {
		return err
	}

	msgs, err := ch.Consume(InvoiceQueue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	sem := make(chan struct{}, runtime.NumCPU()) // buffered channel = concurrency limit

	go func() {
		for msg := range msgs {
			sem <- struct{}{} // blocks here once we're at the cap — this IS the backpressure

			go func(m amqp.Delivery) {
				defer func() { <-sem }()

				var payload InvoiceRequestedPayload
				if err := json.Unmarshal(m.Body, &payload); err != nil {
					log.Printf("bad message, dropping to DLQ: %v", err)
					m.Nack(false, false)
					return
				}

				// context caps how long ONE job is allowed to take —
				// without this, one stuck PDF render blocks a worker slot forever
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				if err := process(ctx, payload); err != nil {
					log.Printf("failed order %s: %v", payload.OrderID, err)
					m.Nack(false, false) // -> DLQ, same as the Nest consumer
					return
				}

				m.Ack(false)
			}(msg)
		}
	}()

	return nil
}
