package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"
	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	DB_DRIVER := viper.GetString("DATABASE.DB_DRIVER")
	DB_NAME := viper.GetString("DATABASE.DB_NAME")

	// DSN untuk SQLite hanya membutuhkan path file database
	dsn := DB_NAME

	db, err := sql.Open(DB_DRIVER, dsn)
	if err != nil {
		log.Fatal("Failed to open connection to database", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database (ping error)", err)
	}

	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT UNIQUE,
		password TEXT
	);`

	_, err = db.Exec(createUserTable)
	if err != nil {
		log.Fatal("Error creating user table:", err)
	} else {
		fmt.Println("Successfully creating user table")
	}

	createFinancialTable := `
	CREATE TABLE IF NOT EXISTS financial_record (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		date DATE NOT NULL,
		type TEXT NOT NULL,
		category TEXT NOT NULL,
		nominal INTEGER NOT NULL,
		description TEXT,
		attachment TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createFinancialTable)
	if err != nil {
		log.Fatal("Error creating financial table:", err)
	} else {
		fmt.Println("Successfully creating financial table")
	}

	log.Println("Database connection successfully")
	return db
}
