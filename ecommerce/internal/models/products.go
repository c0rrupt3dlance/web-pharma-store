package models

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
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
	Product    Product
	Categories []Category
}

type ProductInput struct {
	Product    Product
	Categories []int
}
