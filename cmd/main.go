package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MohamedNagy7/Go-commerce-ops/internal/email"
	"github.com/MohamedNagy7/Go-commerce-ops/internal/invoice"
	"github.com/MohamedNagy7/Go-commerce-ops/internal/pdf"
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

	logoPng, logoErr := os.ReadFile("D:/Go_Lang/go-commerce-ops/internal/assets/favicon.png")
	if logoErr != nil {
		fmt.Println("Error loading logo", logoErr)
		return
	}

	processor := &invoice.Processor{
		GeneratePDF: pdf.GenerateInvoicePDF,
		Email:       email.NewResendProvider(),
		LogoPNG:     logoPng,
	}

	if err := rabbitmq.StartInvoiceConsumer(ch, processor.Process); err != nil {
		fmt.Println("Error starting invoice consumer", err)
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("shutting down")

}
