package entity

import "time"

type ProductPayment struct {
	Id       int
	Name     string
	ImageUrl string
	Price    int
	Qty      int
}

type UserPayment struct {
	UserId        int
	BuyerUsername string
	BuyerName     string
}

type BankPayment struct {
	UserId            int
	BankName          string
	BankAccountName   string
	BankAccountNumber string
}

type Payment struct {
	Id                   string    `json:"id"`
	ProductId            int       `json:"productId"`
	BankAccountId        int       `json:"bankAccountId"`
	PaymentProofImageUrl string    `json:"paymentProofImageUrl"`
	Qty                  int       `json:"quantity"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}
