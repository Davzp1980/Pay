package repository

import (
	"database/sql"
	"fmt"
	"log"
	"pay"
)

//go:generate mockgen -source=createUsers.go -destination=mocks/mock.go
//Поля в postman:
// "name"
// "password"

func CreateAdmin(db *sql.DB, name, hashedPassword string) (string, error) {
	isAdmin := true
	_, err := db.Query("INSERT INTO users (name, password, is_admin) VALUES ($1,$2,$3)", name, hashedPassword, isAdmin)
	if err != nil {
		return "", fmt.Errorf("user %s already exists", name)

	}
	return fmt.Sprintf("Admin %s created", name), nil
	/*return func(w http.ResponseWriter, r *http.Request) {
		var input pay.Input

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		isAdmin := true
		hashedPassword, _ := pay.HashePassword(input.Password)

		_, err = db.Query("INSERT INTO users (name, password, is_admin) VALUES ($1,$2,$3)", input.Name, hashedPassword, isAdmin)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			w.Write([]byte(fmt.Sprintf("User %s already exists", input.Name)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("User %s created", input.Name)))

	}
	*/
}

func CreateUser(db *sql.DB, name, hashedPassword string) (string, error) {
	isAdmin := false
	_, err := db.Query("INSERT INTO users (name, password, is_admin) VALUES ($1,$2,$3)", name, hashedPassword, isAdmin)
	if err != nil {
		return "", fmt.Errorf("user %s already exists", name)

	}
	return fmt.Sprintf("User %s created", name), nil
	/*
		return func(w http.ResponseWriter, r *http.Request) {
			var input pay.Input
			//var user User

			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			isAdmin := false
			hashedPassword, _ := pay.HashePassword(input.Password)

			_, err = db.Query("INSERT INTO users (name, password, is_admin) VALUES ($1,$2,$3)", input.Name, hashedPassword, isAdmin)
			if err != nil {
				log.Println(err)
				w.Write([]byte(fmt.Sprintf("User %s already exists", input.Name)))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("User %s created", input.Name)))
		}
	*/
}

//Поля в postman:
// "name"

func BlockUser(db *sql.DB, name string) (string, error) {
	var user pay.User

	err := db.QueryRow("SELECT name FROM users WHERE name=$1", name).Scan(&user.Name)
	if err != nil {

		return "", err
	}
	expectedName := user.Name
	input_Name := name

	if input_Name == expectedName {
		_, err := db.Exec("UPDATE users SET blocked=$1 WHERE name=$2", true, name)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("user %s does not exists", name)
	}

	return fmt.Sprintf("User %s blocked", name), nil
	/*
		return func(w http.ResponseWriter, r *http.Request) {
			var input pay.Input
			var user pay.User

			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err = db.QueryRow("SELECT name FROM users WHERE name=$1", input.Name).Scan(&user.Name)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			expectedName := user.Name
			inputName := input.Name

			if inputName == expectedName {
				_, err := db.Exec("UPDATE users SET blocked=$1 WHERE name=$2", true, input.Name)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

			} else {
				log.Println("User", inputName, "Does not exists")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("User %s blocked", input.Name)))

		}
	*/
}

func UnBlockUser(db *sql.DB, name string) (string, error) {
	var user pay.User

	err := db.QueryRow("SELECT name FROM users WHERE name=$1", name).Scan(&user.Name)
	if err != nil {

		return "", err
	}
	expectedName := user.Name
	input_Name := name

	if input_Name == expectedName {
		_, err := db.Exec("UPDATE users SET blocked=$1 WHERE name=$2", false, name)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("user %s does not exists", name)
	}

	return fmt.Sprintf("User %s Unblocked", name), nil
	/*
		return func(w http.ResponseWriter, r *http.Request) {
			var input pay.Input
			var user pay.User

			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err = db.QueryRow("SELECT name FROM users WHERE name=$1", input.Name).Scan(&user.Name)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			expectedName := user.Name
			inputName := input.Name

			if inputName == expectedName {
				_, err := db.Exec("UPDATE users SET blocked=$1 WHERE name=$2", false, input.Name)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusBadRequest)
					return

				}

			} else {
				log.Println("User", inputName, "Does not exists")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("User %s unblocked", input.Name)))

		}
	*/
}

func ChangeUserPassword(db *sql.DB, name, hashedPassword string) (string, error) {
	_, err := db.Exec("UPDATE users SET password=$1 WHERE name=$2;", hashedPassword, name)
	if err != nil {
		log.Println("Change password error")
		return "", fmt.Errorf("password change error")
	}

	return fmt.Sprintf("Password %s changed", name), nil
	/*
		return func(w http.ResponseWriter, r *http.Request) {
			var input pay.Input

			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			hashedPassword, _ := pay.HashePassword(input.Password)

			_, err = db.Exec("UPDATE users SET password=$1 WHERE name=$2;", hashedPassword, input.Name)
			if err != nil {
				log.Println("Change password error")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Password %s changed", input.Name)))

		}
	*/

}
