package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pay"
)

//go:generate mockgen -source=createUsers.go -destination=mocks/mock.go
//Поля в postman:
// "name"
// "password"

func CreateAdmin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input

		json.NewDecoder(r.Body).Decode(&input)
		isAdmin := true
		hashedPassword, _ := pay.HashePassword(input.Password)

		_, err := db.Query("INSERT INTO users (name, password, is_admin) VALUES ($1,$2,$3)", input.Name, hashedPassword, isAdmin)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			w.Write([]byte(fmt.Sprintf("User %s already exists", input.Name)))
			return
		}
		w.Write([]byte(fmt.Sprintf("User %s created", input.Name)))

	}

}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		//var user User

		json.NewDecoder(r.Body).Decode(&input)
		isAdmin := false
		hashedPassword, _ := pay.HashePassword(input.Password)

		_, err := db.Query("INSERT INTO users (name, password, is_admin) VALUES ($1,$2,$3)", input.Name, hashedPassword, isAdmin)
		if err != nil {
			log.Println(err)
			w.Write([]byte(fmt.Sprintf("User %s already exists", input.Name)))
			return
		}
		w.Write([]byte(fmt.Sprintf("User %s created", input.Name)))
	}
}

//Поля в postman:
// "name"

func BlockUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		var user pay.User

		json.NewDecoder(r.Body).Decode(&input)

		err := db.QueryRow("SELECT name FROM users WHERE name=$1", input.Name).Scan(&user.Name)
		if err != nil {
			log.Println("Here", err)
		}
		expectedName := user.Name
		inputName := input.Name

		if inputName == expectedName {
			_, err := db.Exec("UPDATE users SET blocked=$1 WHERE name=$2", true, input.Name)
			if err != nil {
				log.Println(err)
			}

		} else {
			log.Println("User", inputName, "Does not exists")
			return
		}
		w.Write([]byte(fmt.Sprintf("User %s blocked", input.Name)))

	}
}

func UnBlockUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input
		var user pay.User

		json.NewDecoder(r.Body).Decode(&input)

		err := db.QueryRow("SELECT name FROM users WHERE name=$1", input.Name).Scan(&user.Name)
		if err != nil {
			log.Println("Here", err)
		}
		expectedName := user.Name
		inputName := input.Name

		if inputName == expectedName {
			_, err := db.Exec("UPDATE users SET blocked=$1 WHERE name=$2", false, input.Name)
			if err != nil {
				log.Println(err)

			}

		} else {
			log.Println("User", inputName, "Does not exists")
			return
		}
		w.Write([]byte(fmt.Sprintf("User %s unblocked", input.Name)))

	}
}

func ChangeUserPassword(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input

		json.NewDecoder(r.Body).Decode(&input)

		hashedPassword, _ := pay.HashePassword(input.Password)

		_, err := db.Exec("UPDATE users SET password=$1 WHERE name=$2;", hashedPassword, input.Name)
		if err != nil {
			log.Println("Change password error")
		}

		w.Write([]byte(fmt.Sprintf("Password %s changed", input.Name)))

	}

}
