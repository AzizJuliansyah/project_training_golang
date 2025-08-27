package controllers

import (
	"database/sql"
	"financial_record/config"
	"financial_record/entities"
	"financial_record/helpers"
	"financial_record/models"
	"financial_record/views"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	db *sql.DB 
}

func NewUserController(db *sql.DB) *UserController {
	return &UserController{db: db}
}

func (controller *UserController) Profile(httpWriter http.ResponseWriter, request *http.Request) {
	templateLayout := "views/user/profile.html"
	
	data := make(map[string]interface{})
	session, _ := config.Store.Get(request, config.SESSION_ID)
	if flashes := session.Flashes("success"); len(flashes) > 0 {
		data["success"] = flashes[0]
		session.Save(request, httpWriter)
	}
	if flashes := session.Flashes("error"); len(flashes) > 0 {
		data["error"] = flashes[0]
		session.Save(request, httpWriter)
	}

	sessionUserId := session.Values["id"].(int)
	userModel := models.NewUserModel(controller.db)
	user, err := userModel.FindUserByID(sessionUserId)
	if err != nil {
		log.Println("Failed to find user data", err)
		data["error"] = "Gagal mendapatkan data user"
	} else {
		data["user"] = user
	}

	if request.Method == http.MethodPost {
		request.ParseMultipartForm(5 << 20)

		password := request.Form.Get("password")
		userInput := entities.User{
			Id: sessionUserId,
			Name: request.Form.Get("name"),
			Email: request.Form.Get("email"),
			Password: password,
		}

		validationResult := helpers.NewValidation().Struct(userInput)
		if validationResult != nil {
			data["validation"] = validationResult
			data["userInput"] = userInput
			views.RenderTemplate(httpWriter, templateLayout, data)
			return
		}

		if password != "" {
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
			userInput.Password = string(hashPassword)
		}

		if file, handler, err := request.FormFile("photo"); err == nil {
			defer file.Close()

			// Validasi ukuran dan tipe
			if handler.Size > 2*1024*1024 {
				data["error"] = "Ukuran file maksimal 2MB"
				data["userInput"] = userInput
				views.RenderTemplate(httpWriter, templateLayout, data)
				return
			}

			ext := strings.ToLower(filepath.Ext(handler.Filename))
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
				data["error"] = "Tipe file harus jpg, jpeg, png, atau webp"
				data["userInput"] = userInput
				views.RenderTemplate(httpWriter, templateLayout, data)
				return
			}

			userModel := models.NewUserModel(controller.db)
			oldPhoto, err := userModel.GetUserPhotoByID(sessionUserId)
			if err != nil {
				log.Println("Failed get user old photo", err)
				data["error"] = "Gagal mendapatkan data Photo."
				data["userInput"] = userInput
				views.RenderTemplate(httpWriter, templateLayout, data)
				return
			}

			// Simpan file ke public/images/news_thumbnail
			filename := fmt.Sprintf("users_photo_profile_%d%s", time.Now().UnixNano(), ext)
			path := filepath.Join("public/images/user_photo_profile", filename)

			// Hapus foto lama jika ada
			if oldPhoto != nil && *oldPhoto != "" {
				oldPath := filepath.Join("public/images/user_photo_profile", *oldPhoto)
				log.Println(oldPath)
				if err := os.Remove(oldPath); err != nil {
					data["error"] = "Gagal menghapus foto lama"
					data["userInput"] = userInput
					views.RenderTemplate(httpWriter, templateLayout, data)
					return
				}
			}

			// Simpan file baru
			out, err := os.Create(path)
			if err != nil {
				data["error"] = "Gagal menyimpan foto"
				data["userInput"] = userInput
				views.RenderTemplate(httpWriter, templateLayout, data)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, file)
			if err != nil {
				data["error"] = "Gagal menyimpan file"
				data["userInput"] = userInput
				views.RenderTemplate(httpWriter, templateLayout, data)
				return
			}

			// Assign ke pointer string
			userInput.Photo = &filename
		} else {
			// Tidak ada file baru → pakai foto lama
			userModel := models.NewUserModel(controller.db)
			oldPhoto, err := userModel.GetUserPhotoByID(sessionUserId)
			if err == nil {
				userInput.Photo = oldPhoto
			}
		}

		err := userModel.UpdateProfile(userInput)
		if err != nil {
			log.Println("Failed to edit profile", err)
			data["error"] = "Gagal mengubah data profile"
		} else {
			session.AddFlash("Berhasil mengubah data profile", "success")
			session.Save(request, httpWriter)
			http.Redirect(httpWriter, request, "/profile", http.StatusSeeOther)
			return
		}

	}

	views.RenderTemplate(httpWriter, templateLayout, data)
}