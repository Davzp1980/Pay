package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"pay"
	"pay/repository"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("BD is not opend")
	}
	defer db.Close()

	repository.MigrateDB(db)

	router := mux.NewRouter()

	router.Use(pay.UserIdentification)

	router.HandleFunc("/sing-in", pay.Login(db)).Methods("POST")
	router.HandleFunc("/sing-up", repository.CreateUser(db)).Methods("POST")
	router.HandleFunc("/logout", pay.Logout).Methods("GET")

	router.HandleFunc("/create-admin", repository.CreateAdmin(db)).Methods("POST")
	router.HandleFunc("/block-user", repository.BlockUser(db)).Methods("POST")
	router.HandleFunc("/unblock-user", repository.UnBlockUser(db)).Methods("POST")

	router.HandleFunc("/change-password", repository.ChangeUserPassword(db)).Methods("POST")

	router.HandleFunc("/block-account", repository.BlockAccount(db)).Methods("POST")
	router.HandleFunc("/unblock-account", repository.UnBlockAccount(db)).Methods("POST")

	router.HandleFunc("/create-account", repository.CreateAccount(db)).Methods("POST")
	router.HandleFunc("/createa-payment", repository.CreatePayment(db)).Methods("POST")
	router.HandleFunc("/replenish-account", repository.ReplenishAccount(db)).Methods("POST")

	router.HandleFunc("/get-account-id", repository.GetAccountsById(db)).Methods("GET")
	router.HandleFunc("/get-account-iban", repository.GetAccountsByIban(db)).Methods("GET")
	router.HandleFunc("/get-account-balance", repository.GetAccountsByBalance(db)).Methods("GET")

	router.HandleFunc("/get-payment-id", repository.GetPaymentsById(db)).Methods("GET")
	router.HandleFunc("/get-payment-date", repository.GetPaymentsDate(db)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))

}
