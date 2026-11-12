package rabbitmq

type InvoiceRequestedPayload struct {
	OrderID     string      `json:"orderId"`
	UserID      string      `json:"userId"`
	Email       string      `json:"email"`
	CustomerName string     `json:"fullName"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"totalAmount"`
	Address     interface{} `json:"address"`
}

type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
