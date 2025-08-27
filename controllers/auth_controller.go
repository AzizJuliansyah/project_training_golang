package controllers

import (
	"database/sql"
	"financial_record/config"
	"financial_record/entities"
	"financial_record/helpers"
	"financial_record/models"
	"financial_record/views"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	db *sql.DB
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{db: db}
}

func (controller *AuthController) Register(httpWriter http.ResponseWriter, request *http.Request) {
	templateLayout := "views/auth/register.html"
	data := make(map[string]interface{})
	session, _ := config.Store.Get(request, config.SESSION_ID)

	data["authInput"] = entities.Register{}
	if request.Method == http.MethodPost {

		request.ParseForm()
		authInput := entities.Register{
			Name:     request.Form.Get("name"),
			Email:    request.Form.Get("email"),
			Password: request.Form.Get("password"),
			ConfirmPassword: request.Form.Get("confirm_password"),
		}

		validationResult := helpers.NewValidation().Struct(authInput)
		if validationResult != nil {
			data["validation"] = validationResult
			data["authInput"] = authInput
			views.RenderTemplate(httpWriter, templateLayout, data)
			return
		}

		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
		authInput.Password = string(hashPassword)

		authModel := models.NewAuthModel(controller.db)
		err := authModel.Register(authInput)
		if err != nil {
			data["error"] = "Gagal register, Terjadi kesalahan server"
			views.RenderTemplate(httpWriter, templateLayout, data)
			return
		}

		session.AddFlash("Berhasil register, silahkan login", "success")
		session.Save(request, httpWriter)

		http.Redirect(httpWriter, request, "/login", http.StatusSeeOther)
		return
	}

	if flashes := session.Flashes("success"); len(flashes) > 0 {
		data["success"] = flashes[0]
		session.Save(request, httpWriter)
	}
	views.RenderTemplate(httpWriter, templateLayout, data)
}

func (controller *AuthController) Login(httpWriter http.ResponseWriter, request *http.Request) {
	templateLayout := "views/auth/login.html"
	data := make(map[string]interface{})
	session, _ := config.Store.Get(request, config.SESSION_ID)

	data["authInput"] = entities.Login{}
	if request.Method == http.MethodPost {

		request.ParseForm()
		authInput := entities.Login{
			Email:    request.Form.Get("email"),
			Password: request.Form.Get("password"),
		}

		validationResult := helpers.NewValidation().Struct(authInput)
		if validationResult != nil {
			data["validation"] = validationResult
			data["authInput"] = authInput
			views.RenderTemplate(httpWriter, templateLayout, data)
			return
		}

		authModel := models.NewAuthModel(controller.db)
		user, err := authModel.FindUserByEmail(authInput.Email)
		if err != nil {
			data["error"] = "Email tidak ditemukan"
			data["authInput"] = authInput
			views.RenderTemplate(httpWriter, templateLayout, data)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authInput.Password)); err != nil {
			data["error"] = "Password salah"
			data["authInput"] = authInput
			views.RenderTemplate(httpWriter, templateLayout, data)
			return
		}

		session, _ := config.Store.Get(request, config.SESSION_ID)
		session.Values["loggedIn"] = true
		session.Values["id"] = user.Id

		session.AddFlash("Berhasil login!", "success")
		session.Save(request, httpWriter)

		http.Redirect(httpWriter, request, "/home", http.StatusSeeOther)
		return
	}

	if flashes := session.Flashes("success"); len(flashes) > 0 {
		data["success"] = flashes[0]
		session.Save(request, httpWriter)
	}
	views.RenderTemplate(httpWriter, templateLayout, data)
}

func Logout(httpWriter http.ResponseWriter, request *http.Request) {
	session, _ := config.Store.Get(request, config.SESSION_ID)

	// kosongkan session
	session.Values = make(map[interface{}]interface{})
	// session.Options.MaxAge = -1

	session.AddFlash("Berhasil logout!", "success")
	session.Save(request, httpWriter)

	http.Redirect(httpWriter, request, "/login", http.StatusSeeOther)
}
