# Go-commerce-ops

An event-driven Go microservice handling background operations for an e-commerce
platform — currently invoice generation and delivery, with OTP delivery and file
storage on the way.

It consumes events published by the main [NestJS backend](#) over RabbitMQ, so
CPU- and I/O-heavy work (PDF rendering, email delivery, future OTP/storage jobs)
runs off the request path.

## Current state

**Implemented**
- Consumes `invoice.requested` events from RabbitMQ
- Generates a branded PDF invoice per order
- Emails the invoice to the customer (via [Resend](https://resend.com)), with
  the PDF and a logo attached inline
- Bounded concurrent processing per message (worker pool sized to CPU cores,
  backed by a semaphore) with per-job timeouts and DLQ routing on failure

**In progress / next**
- Uploading the generated PDF to Cloudflare R2 and persisting a reference to it
  (so invoices are retrievable later, not just delivered once by email)
- OTP generation and delivery (email/SMS) as a second consumer on this service
- Publishing an `invoice.generated` event back to RabbitMQ so the NestJS
  backend can react (e.g. update order status) without polling this service

## Architecture

```
RabbitMQ (invoice.requested)
        │
        ▼
 StartInvoiceConsumer          -- generic message loop, bounded concurrency
        │
        ▼
 invoice.Processor.Process     -- orchestrates one order's invoice
        │
        ├─► BuildInvoiceData    -- maps the queue payload to a domain model
        ├─► pdf.GenerateInvoicePDF
        └─► email.Provider.Send (PDF + logo attached)
```

`internal/invoice` owns the domain model (`InvoiceData`) and orchestration: it
depends on `internal/email` for sending, but nothing in `internal/pdf` or
`internal/email` depends back on it or on RabbitMQ — the PDF renderer and the
email transport don't know a queue exists.

## Project structure

```
cmd/                    entrypoint: wiring only
internal/
  rabbitmq/              connection, topology, generic consumer loop
  invoice/                domain model, payload mapping, orchestration, email template
  pdf/                     PDF rendering
  email/                   email transport (Resend), attachment types
  assets/                  static assets (logo)
```

## Tech stack

- Go
- RabbitMQ (`amqp091-go`)
- [fpdf](https://codeberg.org/go-pdf/fpdf) for PDF generation
- [Resend](https://resend.com) for transactional email

## Running locally

1. Copy `.env.example` to `.env` and fill in:
   ```
   RABBITMQ_URL=amqp://guest:guest@localhost:5672/
   RESEND_API_KEY=your_resend_api_key
   RESEND_FROM_EMAIL=you@yourdomain.com
   ```
2. `go run ./cmd`

The service connects to RabbitMQ, declares its topology (exchange, queue, DLQ)
idempotently, and starts consuming `invoice.requested` messages.

## Reliability

- Messages are acked only after the full pipeline (PDF + email) succeeds;
  failures are nacked straight to a dead-letter queue rather than requeued,
  matching the equivalent NestJS consumer's behavior
- Each job gets a 30s context timeout so one stuck render can't hold a
  worker slot indefinitely