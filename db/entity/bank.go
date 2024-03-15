package entity

type Bank struct {
	Id                string `json:"id"`
	UserId            int    `json:"userId"`
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
}
