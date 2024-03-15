package entity

type Bank struct {
	Id                string `json:"bankAccountId"`
	UserId            int    `json:"-"`
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
}
