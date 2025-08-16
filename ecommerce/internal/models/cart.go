package models

type CartItem struct {
	Id        int     `json:"id"`
	ProductId int     `json:"product_id"`
	Quantity  string  `json:"quantity"`
	Price     float32 `json:"price,omitempty"`
}

type UserCart struct {
	UserId    int        `json:"user_id,omitempty"`
	CartItems []CartItem `json:"cart_items"`
}
