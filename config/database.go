package config

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func InitDB() *sql.DB {
	DB_DRIVER := viper.GetString("DATABASE.DB_DRIVER")
	DB_USER := viper.GetString("DATABASE.DB_USER")
	DB_PASS := viper.GetString("DATABASE.DB_PASS")
	DB_NAME := viper.GetString("DATABASE.DB_NAME")
	DB_HOST := viper.GetString("DATABASE.DB_HOST")
	DB_PORT := viper.GetString("DATABASE.DB_PORT")

	dsn := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?parseTime=true&loc=Asia%2FJakarta"

	db, err := sql.Open(DB_DRIVER, dsn)
	if err != nil {
		log.Fatal("Failed to open connection to database", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database (ping error)", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connection successfully")
	return db
}
