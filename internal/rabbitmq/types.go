package rabbitmq

type InvoiceRequestedPayload struct {
	OrderID     string      `json:"orderId"`
	UserID      string      `json:"userId"`
	Email       string      `json:"email"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"totalAmount"`
	Address     interface{} `json:"address"`
}

type OrderItem struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
