package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// "github.com/MohamedNagy7/Go-commerce-ops/internal/email"
	// "github.com/MohamedNagy7/Go-commerce-ops/internal/pdf"
	"github.com/MohamedNagy7/Go-commerce-ops/internal/rabbitmq"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	//ctx := context.Background()

	rabbitConn, rabbitmqErr := rabbitmq.RabbitMQConnect(os.Getenv("RABBITMQ_URL"))
	if rabbitmqErr != nil {
		fmt.Println("Error connecting to RabbitMQ", rabbitmqErr)
		return
	}
	if rabbitConn == nil {
		fmt.Println("Connection is nil")
		return
	}

	defer rabbitConn.Close()

	ch, channelErr := rabbitConn.Channel()
	if channelErr != nil {
		fmt.Println("Error creating channel", channelErr)
		return
	}
	defer ch.Close()

	if err := rabbitmq.SetupTopology(ch); err != nil {
		fmt.Println("Error setting up topology", err)
		return
	}

	processFunc := func(ctx context.Context, payload rabbitmq.InvoiceRequestedPayload) error {
		fmt.Println("Processing invoice for order:", payload)
		return nil
	}

	testingErr := rabbitmq.StartInvoiceConsumer(ch, processFunc)
	if testingErr != nil {
		fmt.Println("Error starting consumer", testingErr)
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("shutting down")

	// pdfFile, generatePDfErr := pdf.GenerateInvoicePDF()

	// if generatePDfErr != nil {
	// 	fmt.Println("Error generating PDF", generatePDfErr)
	// 	return
	// }

	// emailProvider := email.NewResendProvider()

	// err := emailProvider.Send(ctx, email.Message{
	// 	To:      "mohamed.nagy.khalaf@gmail.com",
	// 	Subject: "Test Email",
	// 	Html:    "<h1>Hello, World!</h1>",
	// }, string(pdfFile))

	// if err != nil {
	// 	fmt.Println("Error sending email", err)
	// 	return
	// }

}
