package models

import (
	"time"
)

type Order struct {
	Id         int         `json:"id"`
	ClientId   int         `json:"client_id"`
	ReceiverId int         `json:"receiver_id"`
	StatusId   int         `json:"status_id"`
	Status     OrderStatus `json:"status"`
	Items      []OrderItem `json:"items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type OrderStatus struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
}

type OrderReceiver struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type OrderItem struct {
	Id       int    `json:"id,omitempty"`
	OrderId  int    `json:"order_id,omitempty"`
	ItemId   int    `json:"item_id"`
	ItemName string `json:"item_name,omitempty"`
	Quantity int    `json:"quantity" binding:"required"`
	Price    int    `json:"price,omitempty"`
}

type OrderItemInput struct {
	ItemId   int `json:"item_id" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}

type OrderInput struct {
	ClientId int              `json:"client_id"`
	Receiver OrderReceiver    `json:"receiver"`
	Items    []OrderItemInput `json:"items"`
}
