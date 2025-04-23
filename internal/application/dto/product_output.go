package dto

type ProductCategoryOutput struct {
	Name string `json:"name"`
}

type ProductImageOutput struct {
	URL string `json:"url"`
}

type ProductOutput struct {
	Name        string                `json:"name"`
	Price       float64               `json:"price"`
	Description string                `json:"description"`
	SKU         string                `json:"sku"`
	Category    ProductCategoryOutput `json:"category"`
	Images      []ProductImageOutput  `json:"images"`
}

type ProductsOutput struct {
	Products    []ProductOutput `json:"products"`
	PageSize    int             `json:"pageSize"`
	CurrentPage int             `json:"currentPage"`
	Total       int             `json:"total"`
}
