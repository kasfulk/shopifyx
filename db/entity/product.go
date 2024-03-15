package entity

type (
	Product struct {
		ID             int      `json:"id"`
		UserID         int      `json:"user_id"`
		Name           string   `json:"name"`
		Price          int      `json:"price"`
		ImageUrl       string   `json:"image_url"`
		Stock          int      `json:"stock"`
		Condition      string   `json:"condition"`
		Tags           []string `json:"tags"`
		IsPurchaseable bool     `json:"is_purchaseable"`
		PurchaseCount  int      `json:"purchase_count"`
	}

	FilterGetProducts struct {
		UserOnly       bool     `json:"userOnly"`
		Limit          int      `json:"limit"`
		Offset         int      `json:"offset"`
		Tags           []string `json:"tags"`
		Condition      string   `json:"condition"`
		ShowEmptyStock bool     `json:"showEmptyStock"`
		MaxPrice       int      `json:"maxPrice"`
		MinPrice       int      `json:"minPrice"`
		SortBy         string   `json:"sortBy"`
		OrderBy        string   `json:"orderBy"`
		Search         string   `json:"search"`
	}
)
