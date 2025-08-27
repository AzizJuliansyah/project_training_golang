package main

import (
	"financial_record/config"
	"financial_record/routes"
	"log"
	"net/http"
)

func main() {
	config.InitViper()

	db := config.InitDB()
	defer db.Close()

	// baca file untuk resource yg ada di folder public (img)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("public/images"))))

	routes.Routes(db)

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}