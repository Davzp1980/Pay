package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"pay"
)

// Поля в postman:
// "name"
// "password"
type outputAccounts struct {
	Id      int    `json:"id"`
	User_id int    `json:"user_id"`
	Iban    string `json:"iban"`
	Balance int    `json:"balance"`
}

func CreateAccount(db *sql.DB, name, iban string) (string, error) {
	var user pay.User
	var account pay.Account
	err := db.QueryRow("SELECT id FROM users WHERE name=$1", name).Scan(&user.ID)
	if err != nil {
		log.Println("User does not exists")
		return "", errors.New("user does not exists")
	}

	err = db.QueryRow("INSERT INTO accounts (user_id, iban) VALUES ($1,$2) RETURNING id", user.ID, iban).Scan(
		&account.ID)
	if err != nil {
		log.Println("account create error")
		return "", errors.New("account create error")
	}
	return fmt.Sprintf("Account %s created", iban), nil
}

func BlockAccount(db *sql.DB, iban string) (string, error) {

	_, err := db.Exec("UPDATE accounts SET blocked=$1 WHERE iban=$2", true, iban)
	if err != nil {

		return "", err
	}
	return fmt.Sprintf("Account %s  is blocked", iban), nil
}

func UnBlockAccount(db *sql.DB, iban string) (string, error) {

	_, err := db.Exec("UPDATE accounts SET blocked=$1 WHERE iban=$2", false, iban)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return fmt.Sprintf("Account %s unblocked", iban), nil
}

func GetAccountsById(db *sql.DB) ([]outputAccounts, error) {
	sortedAccounts := []outputAccounts{}

	rows, err := db.Query("SELECT id, user_id, iban, balance  FROM accounts ORDER BY id")
	if err != nil {
		log.Println("Account selection error")
		return sortedAccounts, err
	}

	for rows.Next() {
		var a outputAccounts

		if err = rows.Scan(&a.Id, &a.User_id, &a.Iban, &a.Balance); err != nil {
			log.Println(err)
		}
		sortedAccounts = append(sortedAccounts, a)
	}

	return sortedAccounts, nil
}

func GetAccountsByIban(db *sql.DB) ([]outputAccounts, error) {

	sortedAccounts := []outputAccounts{}

	rows, err := db.Query("SELECT * FROM accounts ORDER BY iban")
	if err != nil {
		log.Println("Account selection error")
		return sortedAccounts, err
	}

	for rows.Next() {
		var a outputAccounts

		if err = rows.Scan(&a.Id, &a.User_id, &a.Iban, &a.Balance); err != nil {
			log.Println(err)
		}
		sortedAccounts = append(sortedAccounts, a)
	}
	return sortedAccounts, nil
}

func GetAccountsByBalance(db *sql.DB) ([]outputAccounts, error) {

	sortedAccounts := []outputAccounts{}

	rows, err := db.Query("SELECT * FROM accounts ORDER BY balance")
	if err != nil {
		log.Println("Account selection error")
		return sortedAccounts, err
	}

	for rows.Next() {
		var a outputAccounts

		if err = rows.Scan(&a.Id, &a.User_id, &a.Iban, &a.Balance); err != nil {
			log.Println(err)
		}
		sortedAccounts = append(sortedAccounts, a)
	}
	return sortedAccounts, nil
}
