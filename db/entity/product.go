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
	}

	FilterGetProducts struct {
		UserOnly  bool     `json:"userOnly"`
		Limit     int      `json:"limit"`
		Offset    int      `json:"offset"`
		Tags      []string `json:"tags"`
		Condition string   `json:"condition"`
	}
)
