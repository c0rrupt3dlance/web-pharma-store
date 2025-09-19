package models

type Product struct {
	Id          int        `json:"id"`
	Name        string     `json:"name,omitempty" binding:"required"`
	Description string     `json:"description,omitempty" binding:"required"`
	Price       int        `json:"price,omitempty" binding:"required"`
	Media       []MediaUrl `json:"media,omitempty"`
}

type Category struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ProductCategories struct {
	ProductId  int
	CategoryId int
}

type ProductInput struct {
	Product    Product `json:"product"`
	Categories []int   `json:"category_ids,omitempty"`
}

type ProductResponse struct {
	Product    Product     `json:"product"`
	Categories []*Category `json:"categories,omitempty"`
}

type UpdateProductInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Price       *int    `json:"price,omitempty"`
	Categories  []*int  `json:"categories"`
}
