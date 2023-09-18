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
	/*
		return func(w http.ResponseWriter, r *http.Request) {
			var input pay.Input
			var user pay.User
			var account pay.Account

			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err = db.QueryRow("SELECT id FROM users WHERE name=$1", input.Name).Scan(&user.ID)
			if err != nil {
				log.Println("User does not exists")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			i := strconv.Itoa(rand.Intn(1000000000))
			iban := i + input.Name

			err = db.QueryRow("INSERT INTO accounts (user_id, iban) VALUES ($1,$2) RETURNING id", user.ID, iban).Scan(
				&account.ID)
			if err != nil {
				log.Println("Create account error")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			return

		}
	*/
}

func BlockAccount(db *sql.DB, iban string) (string, error) {

	_, err := db.Exec("UPDATE accounts SET blocked=$1 WHERE iban=$2", true, iban)
	if err != nil {

		return "", err
	}
	return fmt.Sprintf("Account %s  is blocked", iban), nil

	/*
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
	*/
}

func UnBlockAccount(db *sql.DB, iban string) (string, error) {

	_, err := db.Exec("UPDATE accounts SET blocked=$1 WHERE iban=$2", false, iban)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return fmt.Sprintf("Account %s unblocked", iban), nil
	/*
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
	*/
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
	/*
		return func(w http.ResponseWriter, r *http.Request) {

			rows, err := db.Query("SELECT id, user_id, iban, balance  FROM accounts ORDER BY id")
			if err != nil {
				log.Panicln("Account selection error")
				w.WriteHeader(http.StatusNotFound)
			}

			sortedAccounts := []outputAccounts{}

			for rows.Next() {
				var a outputAccounts

				if err = rows.Scan(&a.Id, &a.User_id, &a.Iban, &a.Balance); err != nil {
					log.Println(err)
				}
				sortedAccounts = append(sortedAccounts, a)
			}

			json.NewEncoder(w).Encode(sortedAccounts)

		}
	*/
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
	/*
		return func(w http.ResponseWriter, r *http.Request) {

			rows, err := db.Query("SELECT * FROM accounts ORDER BY iban")
			if err != nil {
				log.Panicln("Account selection error")
				w.WriteHeader(http.StatusNotFound)
			}
			sortedAccounts := []outputAccounts{}

			for rows.Next() {
				var a outputAccounts

				if err = rows.Scan(&a.Id, &a.User_id, &a.Iban, &a.Balance); err != nil {
					log.Println(err)
				}
				sortedAccounts = append(sortedAccounts, a)
			}

			json.NewEncoder(w).Encode(sortedAccounts)

		}
	*/
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
	/*
		return func(w http.ResponseWriter, r *http.Request) {

			rows, err := db.Query("SELECT * FROM accounts ORDER BY balance")
			if err != nil {
				log.Panicln("Account selection error")
				w.WriteHeader(http.StatusNotFound)
			}
			sortedAccounts := []outputAccounts{}

			for rows.Next() {
				var a outputAccounts

				if err = rows.Scan(&a.Id, &a.User_id, &a.Iban, &a.Balance); err != nil {
					log.Println(err)
				}
				sortedAccounts = append(sortedAccounts, a)
			}

			json.NewEncoder(w).Encode(sortedAccounts)

		}
	*/
}
