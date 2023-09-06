package pay

import "time"

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	IsAdmin      bool   `json:"is_admin"`
	Blocked      bool   `json:"blocked"`
}

type Account struct {
	ID      int    `json:"id"`
	UserId  int    `json:"user_id"`
	Iban    string `json:"iban"`
	Balance int    `json:"balance"`
	Blocked bool   `json:"blocked"`
}

type Payment struct {
	ID            int       `json:"id"`
	UserId        int       `json:"user_id"`
	Reciever      string    `json:"reciever"`
	RecieverIban  string    `json:"reciever_iban"`
	Payer         string    `json:"payer"`
	PayerIban     string    `json:"payer_iban"`
	AmountPayment int       `json:"amount_payment"`
	Date          time.Time `json:"date"`
}

type Input struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Sort     string `json:"sort"`
	Iban     string `json:"iban"`
}

type InputPayment struct {
	PayerName     string `json:"payer_name"`
	PayerIban     string `json:"payer_iban"`
	ReceiverName  string `json:"receiver_name"`
	ReceiverIban  string `json:"receiver_iban"`
	AmountPayment int    `json:"amount_payment"`
}
