package entity

type Product struct {
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
