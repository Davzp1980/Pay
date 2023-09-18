package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pay"
	"pay/repository"
)

func CreatePayment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.InputPayment

		json.NewDecoder(r.Body).Decode(&input)

		response, err := repository.CreatePayment(db, input)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusOK)
	}
}

func GetPaymentsById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response, err := repository.GetPaymentsById(db)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusOK)
	}
}

func GetPaymentsDate(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response, err := repository.GetPaymentsDate(db)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusOK)
	}
}

func ReplenishAccount(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input

		json.NewDecoder(r.Body).Decode(&input)

		response, err := repository.ReplenishAccount(db, input.Name, input.Iban, input.AmountReplenish)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}
}
