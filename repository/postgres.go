package repository

import "database/sql"

func MigrateDB(db *sql.DB) error {

	CreateTablesQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id serial PRIMARY KEY,
			name VARCHAR (25) NOT NULL UNIQUE,
			password VARCHAR NOT NULL UNIQUE,
			is_admin boolean,
			blocked boolean DEFAULT false
		);
		CREATE TABLE IF NOT EXISTS accounts (
			id serial PRIMARY KEY,
			user_id serial NOT NULL,
			iban VARCHAR (34) NOT NULL UNIQUE,
			balance integer DEFAULT 0,
			blocked boolean

		);
		CREATE TABLE IF NOT EXISTS payments (
			id serial PRIMARY KEY,
			user_id serial NOT NULL,
			reciever VARCHAR,
			reciever_iban VARCHAR,
			payer VARCHAR,
			payer_iban VARCHAR,
			amount_payment serial,
			date TIMESTAMP WITH TIME ZONE
		);
	`

	_, err := db.Exec(CreateTablesQuery)
	return err
}
