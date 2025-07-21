package models

import "time"

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name,omitempty" binding:"required"`
	Description string  `json:"description,omitempty" binding:"required"`
	Price       float32 `json:"price,omitempty" binding:"required"`
}

type Category struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ProductCategories struct {
	ProductId  int
	CategoryId int
}

type ProductResponse struct {
	Product    Product     `json:"product"`
	Categories []*Category `json:"categories,omitempty"`
}

type ProductInput struct {
	Product    Product `json:"product"`
	Categories []int   `json:"category_ids,omitempty"`
}

type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float32 `json:"price,omitempty"`
	Categories  []*int   `json:"category_ids"`
}

type Media struct {
	Id          int       `json:"id"`
	Bucket      string    `json:"bucket"`
	ObjectKey   string    `json:"object_key"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductMedia struct {
	Id        int
	ProductId int
	MediaId   int
}
