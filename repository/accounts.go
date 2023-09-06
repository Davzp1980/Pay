package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"pay"
	"strconv"
)

//Поля в postman:
// "name"
// "password"

func CreateAccount(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		var user pay.User
		var account pay.Account

		json.NewDecoder(r.Body).Decode(&input)

		err := db.QueryRow("SELECT id FROM users WHERE name=$1", input.Name).Scan(&user.ID)
		if err != nil {
			log.Println("User does not exists")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		i := strconv.Itoa(rand.Intn(1000000000))
		iban := i + input.Name
		fmt.Println(iban)

		err = db.QueryRow("INSERT INTO accounts (user_id, iban) VALUES ($1,$2) RETURNING id", user.ID, iban).Scan(
			&account.ID)
		if err != nil {
			log.Println("Create account error")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}

func BlockAccount(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		var account pay.Account

		json.NewDecoder(r.Body).Decode(&input)

		err := db.QueryRow("SELECT iban, blocked FROM accounts WHERE iban=$1", input.Iban).Scan(&account.Iban, &account.Blocked)
		if err != nil {
			log.Println("Here", err)
		}
		expectedIban := account.Iban
		inputIban := input.Iban

		if inputIban == expectedIban {
			_, err := db.Exec("UPDATE accounts SET blocked=$1 WHERE iban=$2", true, input.Iban)
			if err != nil {
				log.Println(err)
			}

		} else {
			log.Println("Account", inputIban, "Does not exists")
			return
		}
		w.Write([]byte(fmt.Sprintf("Account %s blocked", input.Name)))

	}
}

func UnBlockAccount(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		var account pay.Account

		json.NewDecoder(r.Body).Decode(&input)

		err := db.QueryRow("SELECT iban, blocked FROM accounts WHERE iban=$1", input.Iban).Scan(&account.Iban, &account.Blocked)
		if err != nil {
			log.Println("Here", err)
		}
		expectedIban := account.Iban
		inputIban := input.Iban

		if inputIban == expectedIban {
			_, err := db.Exec("UPDATE accounts SET blocked=$1 WHERE iban=$2", false, input.Iban)
			if err != nil {
				log.Println(err)
			}

		} else {
			log.Println("Account", inputIban, "Does not exists")
			return
		}
		w.Write([]byte(fmt.Sprintf("Account %s unblocked", input.Name)))

	}
}

func GetAccountsById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT * FROM accounts ORDER BY id")
		if err != nil {
			log.Panicln("Account selection error")
			w.WriteHeader(http.StatusNotFound)
		}

		sortedAccounts := []pay.Account{}

		for rows.Next() {
			var a pay.Account

			if err = rows.Scan(&a.ID, &a.UserId, &a.Iban, &a.Balance); err != nil {
				log.Println(err)
			}
			sortedAccounts = append(sortedAccounts, a)
		}

		json.NewEncoder(w).Encode(sortedAccounts)

	}
}

func GetAccountsByIban(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT * FROM accounts ORDER BY iban")
		if err != nil {
			log.Panicln("Account selection error")
			w.WriteHeader(http.StatusNotFound)
		}
		sortedAccounts := []pay.Account{}

		for rows.Next() {
			var a pay.Account

			if err = rows.Scan(&a.ID, &a.UserId, &a.Iban, &a.Balance); err != nil {
				log.Println(err)
			}
			sortedAccounts = append(sortedAccounts, a)
		}

		json.NewEncoder(w).Encode(sortedAccounts)

	}
}

func GetAccountsByBalance(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT * FROM accounts ORDER BY balance")
		if err != nil {
			log.Panicln("Account selection error")
			w.WriteHeader(http.StatusNotFound)
		}
		sortedAccounts := []pay.Account{}

		for rows.Next() {
			var a pay.Account

			if err = rows.Scan(&a.ID, &a.UserId, &a.Iban, &a.Balance); err != nil {
				log.Println(err)
			}
			sortedAccounts = append(sortedAccounts, a)
		}

		json.NewEncoder(w).Encode(sortedAccounts)

	}
}
