package repository

import (
	"database/sql"
	"fmt"
	"log"
	"pay"
	"time"
)

/*
			поля при создании платежа:
		"payer_name":"alex",
	    "payer_iban":"012588alex",
	    "amount_payment":1300,

	    "receiver_name":"ira",
	    "receiver_iban":"27887ira"
*/
func CreatePayment(db *sql.DB, input pay.InputPayment) (string, error) {
	var payer pay.User
	var receiverAccount pay.Account
	var payerAccount pay.Account
	var payment pay.Payment

	//по имени отправителя получаем id
	err := db.QueryRow("SELECT id, name FROM users WHERE name=$1", input.PayerName).Scan(&payer.ID, &payer.Name)
	if err != nil {
		log.Println("User does not exists")
		return "", err

	}
	// по номеу счета (iban) получаем id, Iban, Balance получателя
	err = db.QueryRow("SELECT id, user_id, iban, balance, blocked  FROM accounts WHERE iban=$1", input.ReceiverIban).Scan(
		&receiverAccount.ID, &receiverAccount.UserId, &receiverAccount.Iban, &receiverAccount.Balance, &receiverAccount.Blocked)
	if err != nil {
		log.Println("Account does not exists")
		return "", err
	}
	if receiverAccount.Blocked {
		log.Println("Reciever account blocked")
		return "", err
	}
	// создаем платеж
	err = db.QueryRow("INSERT INTO payments (user_id, reciever, reciever_iban, payer, payer_iban, amount_payment, date) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
		receiverAccount.UserId, input.ReceiverName, input.ReceiverIban, input.PayerName, input.PayerIban, input.AmountPayment, time.Now()).Scan(
		&payment.ID)
	if err != nil {
		log.Println("Create payment error")
		return "", err
	}
	// проверяем достаточно ли денег на счете отправителя и снимаем сумму платежа со счета
	err = db.QueryRow("SELECT balance, blocked FROM accounts WHERE iban=$1", input.PayerIban).Scan(&payerAccount.Balance, &payerAccount.Blocked)
	if err != nil {
		log.Println("Wrong payer balance")
		return "", err
	}
	if payerAccount.Blocked {
		log.Println("Payer account blocked")
		return "", err
	}

	if payerAccount.Balance < input.AmountPayment {
		log.Println("Not enough money in the account")
		return "", err
	}
	payerBalance := payerAccount.Balance - input.AmountPayment

	_, err = db.Exec("UPDATE accounts SET balance=$2 WHERE iban=$1", input.PayerIban, payerBalance)
	if err != nil {
		log.Println("Add to balance error")
		return "", err
	}
	// изменяем баланс получателя в соответствии с указынным номером счета (iban) и суммой платежа
	balance := receiverAccount.Balance + input.AmountPayment

	_, err = db.Exec("UPDATE accounts SET balance=$2 WHERE iban=$1", input.ReceiverIban, balance)
	if err != nil {
		log.Println("Add to balance error")
		return "", err
	}
	return fmt.Sprintf("Payment payment was made %v UAH", input.AmountPayment), nil

	/*
		return func(w http.ResponseWriter, r *http.Request) {
			var input pay.InputPayment
			var payer pay.User
			var receiverAccount pay.Account
			var payerAccount pay.Account
			var payment pay.Payment

			json.NewDecoder(r.Body).Decode(&input)
			//по имени отправителя получаем id
			err := db.QueryRow("SELECT id, name FROM users WHERE name=$1", input.PayerName).Scan(&payer.ID, &payer.Name)
			if err != nil {
				log.Println("User does not exists")
				w.WriteHeader(http.StatusForbidden)
			}
			// по номеу счета (iban) получаем id, Iban, Balance получателя
			err = db.QueryRow("SELECT id, user_id, iban, balance, blocked  FROM accounts WHERE iban=$1", input.ReceiverIban).Scan(
				&receiverAccount.ID, &receiverAccount.UserId, &receiverAccount.Iban, &receiverAccount.Balance, &receiverAccount.Blocked)
			if err != nil {
				log.Println("Account does not exists")
				w.WriteHeader(http.StatusForbidden)
			}
			if receiverAccount.Blocked {
				log.Println("Reciever account blocked")
				return
			}
			// создаем платеж
			err = db.QueryRow("INSERT INTO payments (user_id, reciever, reciever_iban, payer, payer_iban, amount_payment, date) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
				receiverAccount.UserId, input.ReceiverName, input.ReceiverIban, input.PayerName, input.PayerIban, input.AmountPayment, time.Now()).Scan(
				&payment.ID)
			if err != nil {
				log.Println("Create payment error")
				w.WriteHeader(http.StatusForbidden)
			}
			// проверяем достаточно ли денег на счете отправителя и снимаем сумму платежа со счета
			err = db.QueryRow("SELECT balance, blocked FROM accounts WHERE iban=$1", input.PayerIban).Scan(&payerAccount.Balance, &payerAccount.Blocked)
			if err != nil {
				log.Println("Wrong payer balance")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if payerAccount.Blocked {
				log.Println("Payer account blocked")
				return
			}

			if payerAccount.Balance < input.AmountPayment {
				log.Println("Not enough money in the account")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			payerBalance := payerAccount.Balance - input.AmountPayment

			_, err = db.Exec("UPDATE accounts SET balance=$2 WHERE iban=$1", input.PayerIban, payerBalance)
			if err != nil {
				log.Println("Add to balance error")
				w.WriteHeader(http.StatusForbidden)
			}
			// изменяем баланс получателя в соответствии с указынным номером счета (iban) и суммой платежа
			balance := receiverAccount.Balance + input.AmountPayment

			_, err = db.Exec("UPDATE accounts SET balance=$2 WHERE iban=$1", input.ReceiverIban, balance)
			if err != nil {
				log.Println("Add to balance error")
				w.WriteHeader(http.StatusForbidden)
			}
			w.Write([]byte(fmt.Sprintf("Payment payment was made %v UAH", input.AmountPayment)))

		}

	*/
}

// получение платежей по имени получателя с сортировкой
func GetPaymentsById(db *sql.DB) ([]pay.Payment, error) {

	sortedPayments := []pay.Payment{}

	rows, err := db.Query("SELECT * FROM payments ORDER BY id")
	if err != nil {
		log.Println("Account selection error")
		return sortedPayments, err
	}

	for rows.Next() {
		var p pay.Payment

		if err = rows.Scan(&p.ID, &p.UserId, &p.Reciever, &p.RecieverIban, &p.Payer, &p.PayerIban, &p.AmountPayment, &p.Date); err != nil {
			log.Println(err)
		}
		sortedPayments = append(sortedPayments, p)
	}

	return sortedPayments, nil
	/*
		return func(w http.ResponseWriter, r *http.Request) {

			rows, err := db.Query("SELECT * FROM payments ORDER BY id")
			if err != nil {
				log.Panicln("Account selection error")
				w.WriteHeader(http.StatusNotFound)
			}

			sortedPayments := []pay.Payment{}

			for rows.Next() {
				var p pay.Payment

				if err = rows.Scan(&p.ID, &p.UserId, &p.Reciever, &p.RecieverIban, &p.Payer, &p.PayerIban, &p.AmountPayment, &p.Date); err != nil {
					log.Println(err)
				}
				sortedPayments = append(sortedPayments, p)
			}

			json.NewEncoder(w).Encode(sortedPayments)

		}
	*/
}

func GetPaymentsDate(db *sql.DB) ([]pay.Payment, error) {

	sortedPayments := []pay.Payment{}

	rows, err := db.Query("SELECT * FROM payments ORDER BY date DESC")
	if err != nil {
		log.Println("Account selection error")
		return sortedPayments, err
	}

	for rows.Next() {
		var p pay.Payment

		if err = rows.Scan(&p.ID, &p.UserId, &p.Reciever, &p.RecieverIban, &p.Payer, &p.PayerIban, &p.AmountPayment, &p.Date); err != nil {
			log.Println(err)
		}
		sortedPayments = append(sortedPayments, p)
	}
	return sortedPayments, nil
	/*
		return func(w http.ResponseWriter, r *http.Request) {

			rows, err := db.Query("SELECT * FROM payments ORDER BY date DESC")
			if err != nil {
				log.Panicln("Account selection error")
				w.WriteHeader(http.StatusNotFound)
			}

			sortedPayments := []pay.Payment{}

			for rows.Next() {
				var p pay.Payment

				if err = rows.Scan(&p.ID, &p.UserId, &p.Reciever, &p.RecieverIban, &p.Payer, &p.PayerIban, &p.AmountPayment, &p.Date); err != nil {
					log.Println(err)
				}
				sortedPayments = append(sortedPayments, p)
			}

			json.NewEncoder(w).Encode(sortedPayments)

		}
	*/
}

func ReplenishAccount(db *sql.DB, name, iban string, amountReplenish int) (string, error) {
	var user pay.User
	var amountAccount pay.Account

	err := db.QueryRow("SELECT id FROM users WHERE name=$1", name).Scan(&user.ID)
	if err != nil {
		log.Println("User does not exists")
		return "User does not exists", err
	}

	err = db.QueryRow("SELECT balance FROM accounts WHERE iban=$1", iban).Scan(&amountAccount.Balance)
	if err != nil {
		log.Println("Account does not exists")
		return "Account does not exists", err
	}
	balance := amountReplenish + amountAccount.Balance

	_, err = db.Exec("UPDATE accounts SET balance=$1 WHERE iban=$2", balance, iban)
	if err != nil {
		log.Println("UPDATE account error")
		return "UPDATE account error", err
	}
	return fmt.Sprintf("Account was replenished for %v. Amount in the account %v", amountReplenish, balance), nil

	/*
		return func(w http.ResponseWriter, r *http.Request) {
			var input pay.Input
			var user pay.User
			var amountAccount int

			json.NewDecoder(r.Body).Decode(&input)
			//balance:=input.AmountReplenish+
			err := db.QueryRow("SELECT id FROM users WHERE name=$1", input.Name).Scan(&user.ID)
			if err != nil {
				log.Println("User does not exists")
				w.WriteHeader(http.StatusBadRequest)
			}

			err = db.QueryRow("SELECT balance FROM accounts WHERE user_id=$1", user.ID).Scan(&amountAccount)
			if err != nil {
				log.Println("Account does not exists")
				w.WriteHeader(http.StatusInternalServerError)
			}
			balance := input.AmountReplenish + amountAccount

			_, err = db.Exec("UPDATE accounts SET balance=$1", balance)
			if err != nil {
				log.Println("UPDATE account error")
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write([]byte(fmt.Sprintf("Account was replenished for %v\n", input.AmountReplenish)))
			w.Write([]byte(fmt.Sprintf("Amount in the account %v", balance)))
		}
	*/

}
