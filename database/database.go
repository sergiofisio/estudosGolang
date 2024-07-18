package database

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {
    createTablesSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		document VARCHAR(255) NOT NULL UNIQUE,
        username VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL
    );

    CREATE TABLE IF NOT EXISTS addresses (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        street VARCHAR(255) NOT NULL,
		number INTEGER NOT NULL,
		complement VARCHAR(255),
		neighborhood VARCHAR(255) NOT NULL,
        city VARCHAR(255) NOT NULL,
        state VARCHAR(255) NOT NULL,
		country VARCHAR(255) NOT NULL,
        zip_code VARCHAR(255) NOT NULL
    );

    CREATE TABLE IF NOT EXISTS phones (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
		phone_type VARCHAR(255) NOT NULL,
		country_code VARCHAR(255) NOT NULL,
		area_code VARCHAR(255) NOT NULL,
        phone_number VARCHAR(255) NOT NULL
    );
	
    CREATE TABLE IF NOT EXISTS purchases (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        product_name VARCHAR(255) NOT NULL,
        amount DECIMAL(10, 2) NOT NULL,
        purchase_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    `

    _, err := db.Exec(createTablesSQL)
    if err != nil {
        log.Fatal("Falha ao criar tabelas: ", err)
    }
}