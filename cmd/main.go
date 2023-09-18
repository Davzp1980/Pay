package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"pay"

	"pay/handler"
	"pay/repository"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	//password := os.Getenv("DB_PASSWORD") fmt.Sprintf("host=go_db user=postgres password=%v dbname=mydb sslmode=disable", password
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("BD is not opend")
	}
	defer db.Close()

	repository.MigrateDB(db)

	router := mux.NewRouter()

	router.Use(pay.UserIdentification)

	router.HandleFunc("/sign-in", pay.Login(db)).Methods("POST")
	router.HandleFunc("/sign-up", handler.CreateUser(db)).Methods("POST")
	router.HandleFunc("/logout", pay.Logout).Methods("GET")

	router.HandleFunc("/create-admin", handler.CreateAdmin(db)).Methods("POST")
	router.HandleFunc("/block-user", handler.BlockUser(db)).Methods("POST")
	router.HandleFunc("/unblock-user", handler.UnBlockUser(db)).Methods("POST")

	router.HandleFunc("/change-password", handler.ChangeUserPassword(db)).Methods("POST")

	router.HandleFunc("/block-account", handler.BlockAccount(db)).Methods("POST")
	router.HandleFunc("/unblock-account", handler.UnBlockAccount(db)).Methods("POST")

	router.HandleFunc("/create-account", handler.CreateAccount(db)).Methods("POST")
	router.HandleFunc("/createa-payment", handler.CreatePayment(db)).Methods("POST")
	router.HandleFunc("/replenish-account", handler.ReplenishAccount(db)).Methods("POST")

	router.HandleFunc("/get-account-id", handler.GetAccountsById(db)).Methods("GET")
	router.HandleFunc("/get-account-iban", handler.GetAccountsByIban(db)).Methods("GET")
	router.HandleFunc("/get-account-balance", handler.GetAccountsByBalance(db)).Methods("GET")

	router.HandleFunc("/get-payment-id", handler.GetPaymentsById(db)).Methods("GET")
	router.HandleFunc("/get-payment-date", handler.GetPaymentsDate(db)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))

}
