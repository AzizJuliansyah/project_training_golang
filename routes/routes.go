package routes

import (
	"database/sql"
	"financial_record/controllers"
	"financial_record/helpers"
	"net/http"
)

func Routes(db *sql.DB) {
	// auth routes
	authController := controllers.NewAuthController(db)
	http.HandleFunc("/", helpers.GuestOnly(authController.Login))
	http.HandleFunc("/register", helpers.GuestOnly(authController.Register))
	http.HandleFunc("/login", helpers.GuestOnly(authController.Login))
	http.HandleFunc("/logout", helpers.AuthOnly(controllers.Logout))

	userController := controllers.NewUserController(db)
	http.HandleFunc("/profile", helpers.AuthOnly(userController.Profile))

	financialController := controllers.NewFinancialController(db)
	http.HandleFunc("/home", helpers.AuthOnly(financialController.Home))
	http.HandleFunc("/financial/add_financial_record", helpers.AuthOnly(financialController.AddFinancialRecord))
	http.HandleFunc("/financial/edit_financial_record", helpers.AuthOnly(financialController.EditFinancialRecord))
	http.HandleFunc("/financial/delete_financial_record", helpers.AuthOnly(financialController.DeleteFinancial))
	http.HandleFunc("/financial/download_financial_record", helpers.AuthOnly(financialController.DownloadFinancialRecord))


}