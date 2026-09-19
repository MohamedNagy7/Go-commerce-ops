package invoice

import "time"

type LineItem struct {
	ProductName string
	Quantity    int
	Price       float64
}

type InvoiceData struct {
	OrderID       string
	CustomerEmail string
	Items         []LineItem
	Subtotal      float64
	Total         float64
	IssuedAt      time.Time
}
