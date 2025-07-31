package models

type CartItem struct {
	ProductId   int     `json:"product_id"`
	Description string  `json:"description"`
	Quantity    string  `json:"quantity"`
	Price       float32 `json:"price"`
}
