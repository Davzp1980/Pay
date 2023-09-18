package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pay"
	"pay/repository"
)

func CreateAdmin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		var response string
		json.NewDecoder(r.Body).Decode(&input)

		hashedPassword, _ := pay.HashePassword(input.Password)

		response, err := repository.CreateAdmin(db, input.Name, hashedPassword)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		var response string
		json.NewDecoder(r.Body).Decode(&input)

		hashedPassword, _ := pay.HashePassword(input.Password)
		response, err := repository.CreateUser(db, input.Name, hashedPassword)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}

func BlockUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		var response string
		json.NewDecoder(r.Body).Decode(&input)

		response, err := repository.BlockUser(db, input.Name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}

func UnBlockUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		var response string
		json.NewDecoder(r.Body).Decode(&input)

		response, err := repository.UnBlockUser(db, input.Name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}

func ChangeUserPassword(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		var response string
		json.NewDecoder(r.Body).Decode(&input)

		hashedPassword, _ := pay.HashePassword(input.Password)

		response, err := repository.ChangeUserPassword(db, input.Name, hashedPassword)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}
