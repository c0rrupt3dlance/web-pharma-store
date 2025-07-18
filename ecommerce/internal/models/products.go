package models

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

type ProductImage struct {
	Id        int    `json:"id"`
	ProductId int    `json:"product_id"`
	ImageUrl  string `json:"image_url"`
	AltText   string `json:"alt_text"`
	IsMain    bool   `json:"is_main"`
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
	Categories  []*int   `json:"categories"`
}
