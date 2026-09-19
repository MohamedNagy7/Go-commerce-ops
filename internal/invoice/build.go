package invoice

import (
	"time"

	"github.com/MohamedNagy7/Go-commerce-ops/internal/rabbitmq"
)

func BuildInvoiceData(payload rabbitmq.InvoiceRequestedPayload) InvoiceData {
	items := make([]LineItem, len(payload.Items))
	var subtotal float64

	for i, item := range payload.Items {
		items[i] = LineItem{ProductName: item.ProductName, Quantity: item.Quantity, Price: item.Price}
		subtotal += item.Price
	}

	return InvoiceData{
		OrderID:       payload.OrderID,
		CustomerEmail: payload.Email,
		Items:         items,
		Subtotal:      subtotal,
		Total:         payload.TotalAmount,
		IssuedAt:      time.Now(),
	}
}
