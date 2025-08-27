package controllers

import (
	"database/sql"
	"financial_record/config"
	"financial_record/entities"
	"financial_record/helpers"
	"financial_record/models"
	"financial_record/views"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FinancialController struct {
	db *sql.DB
}

func NewFinancialController(db *sql.DB) *FinancialController {
	return &FinancialController{db: db}
}

func formatIDR(n int64) string {
	str := fmt.Sprintf("%d", n)
	var result []string
	for len(str) > 3 {
		result = append([]string{str[len(str)-3:]}, result...)
		str = str[:len(str)-3]
	}
	if len(str) > 0 {
		result = append([]string{str}, result...)
	}
	return strings.Join(result, ".") + ",00"
}

func (controller *FinancialController) Home(httpWriter http.ResponseWriter, request *http.Request) {
	templateLayout := "views/financial/home.html"

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

	currentDate := time.Now()
	var months []string
	for i := 0; i < 12; i++ {
		previousMonth := currentDate.AddDate(0, -i, 0)
		months = append(months, previousMonth.Format("January 2006"))
	}
	
	data["months"] = months
	selectedMonth := request.URL.Query().Get("selected_month")
	log.Println("ctr", selectedMonth)
	if selectedMonth == "" {
		selectedMonth = currentDate.Format("January 2006")
	}
	data["selectedMonth"] = selectedMonth

	pemasukanOnly := request.URL.Query().Get("pemasukanOnly") == "true"
	data["pemasukanOnly"] = pemasukanOnly
	pengeluaranOnly := request.URL.Query().Get("pengeluaranOnly") == "true"
	data["pengeluaranOnly"] = pengeluaranOnly

	sessionUserId := session.Values["id"].(int)
	financialModel := models.NewFinancialModel(controller.db)
	total_pemasukan, total_pengeluaran, err := financialModel.GetFinancialTotalNominal(sessionUserId, selectedMonth, pemasukanOnly, pengeluaranOnly)
	if err != nil {
		log.Println("Failed to get financial total nominal data", err)
		data["error"] = "Gagal mendapatkan total data keuangan"
	} else {
		data["total_pemasukan"] = total_pemasukan
		data["total_pengeluaran"] = total_pengeluaran
	}

	financials, err := financialModel.FindAllFinancial(sessionUserId, selectedMonth, pemasukanOnly, pengeluaranOnly)
	if err != nil {
		log.Println("Failed to get all financials data", err)
		data["error"] = "Gagal mendapatkan semua data keuangan"
	} else {
		data["financials"] = financials
	}

	funcMap := template.FuncMap{
		"formatIDR": formatIDR,
		"indexNo": func(a, b int) int { return a + b },
	}
	tmpl, _ := template.New(filepath.Base(templateLayout)).Funcs(funcMap).ParseFiles(templateLayout)
	_ = tmpl.Execute(httpWriter, data)
}

func (controller *FinancialController) AddFinancialRecord(httpWriter http.ResponseWriter, request *http.Request) {
	templateLayout := "views/financial/create.html"
	
	data := make(map[string]interface{})
	session, _ := config.Store.Get(request, config.SESSION_ID)

	if request.Method == http.MethodPost {
		request.ParseForm()

		dateStr := request.Form.Get("date")
		date, _ := time.Parse("2006-01-02", dateStr)

		nominalStr := request.Form.Get("nominal")
		NominalInt64, _ := strconv.ParseInt(nominalStr, 10, 64)

		var attachment *string
		if attachmentValue := request.Form.Get("attachment"); attachmentValue != "" {
			attachment = &attachmentValue
		}

		var description *string
		if descriptionValue := request.Form.Get("description"); descriptionValue != "" {
			description = &descriptionValue
		}

		sessionUserId := session.Values["id"].(int)
		financial := entities.AddFinancial{
			UserId:		 sessionUserId,
			Date:        date,
			Type:        request.Form.Get("type"),
			Nominal: 	 NominalInt64,
			Category:    request.Form.Get("category"),
			Description: description,
			Attachment:  attachment,
		}
		log.Println(financial)

		validationResult := helpers.NewValidation().Struct(financial)
		if validationResult != nil {
			data["validation"] = validationResult
			views.RenderTemplate(httpWriter, templateLayout, data)
			return
		}

		
		financialModel := models.NewFinancialModel(controller.db)
		err := financialModel.AddFinancialRecord(financial)
		if err != nil {
			log.Println("Failed to add financial record", err)
			data["error"] = "Gagal menambahkan data keuangan"
		} else {
			session.AddFlash("Berhasil menambahkan data keuangan", "success")
			session.Save(request, httpWriter)
			http.Redirect(httpWriter, request, "/home", http.StatusSeeOther)
			return
		}
	}

	views.RenderTemplate(httpWriter, templateLayout, data)
}

func (controller *FinancialController) EditFinancialRecord(httpWriter http.ResponseWriter, request *http.Request) {
	templateLayout := "views/financial/edit.html"
	
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

	id := request.URL.Query().Get("id")
	Int64, err := strconv.ParseInt(id, 10, 64)
	if id == "" || err != nil {
		session.AddFlash("Gagal mendapatkan data keuangan, ID kosong!", "error")
		session.Save(request, httpWriter)
	
		http.Redirect(httpWriter, request, "/home", http.StatusSeeOther)
		return
	}

	financialModels := models.NewFinancialModel(controller.db)
	finacial, err := financialModels.FindFinancialByID(Int64)
	if err != nil {
		log.Println("Failed to find financial data", err)
		data["error"] = "Gagal menadapatkan data keuangan"
	} else {
		data["financial"] = finacial
	}

	if request.Method == http.MethodPost {
		request.ParseForm()

		dateStr := request.Form.Get("date")
		date, _ := time.Parse("2006-01-02", dateStr)

		nominalStr := request.Form.Get("nominal")
		NominalInt64, _ := strconv.ParseInt(nominalStr, 10, 64)

		var attachment *string
		if attachmentValue := request.Form.Get("attachment"); attachmentValue != "" {
			attachment = &attachmentValue
		}

		var description *string
		if descriptionValue := request.Form.Get("description"); descriptionValue != "" {
			description = &descriptionValue
		}

		financial := entities.AddFinancial{
			Id: 		 Int64,
			Date:        date,
			Type:        request.Form.Get("type"),
			Nominal: 	 NominalInt64,
			Category:    request.Form.Get("category"),
			Description: description,
			Attachment:  attachment,
		}
		log.Println(financial)

		validationResult := helpers.NewValidation().Struct(financial)
		if validationResult != nil {
			data["validation"] = validationResult
			views.RenderTemplate(httpWriter, templateLayout, data)
			return
		}

		
		financialModel := models.NewFinancialModel(controller.db)
		err := financialModel.EditFinancialRecord(financial)
		if err != nil {
			log.Println("Failed to edit financial record", err)
			data["error"] = "Gagal mengubah data keuangan"
		} else {
			session.AddFlash("Berhasil mengubah data keuangan", "success")
			session.Save(request, httpWriter)
			http.Redirect(httpWriter, request, "/financial/edit_financial_record?id=" + id, http.StatusSeeOther)
			return
		}
	}

	views.RenderTemplate(httpWriter, templateLayout, data)
}

func (controller *FinancialController) DeleteFinancial(httpWriter http.ResponseWriter, request *http.Request) {
	session, _ := config.Store.Get(request, config.SESSION_ID)


	id := request.URL.Query().Get("id")
	int64Id, err := strconv.ParseInt(id, 10, 64)
	if id == "" || err != nil {
		session.AddFlash("Gagal mendapatkan data keuangan, ID kosong!", "error")
		session.Save(request, httpWriter)
	
		http.Redirect(httpWriter, request, "/home", http.StatusSeeOther)
		return
	}

	financialModels := models.NewFinancialModel(controller.db)
	err = financialModels.DeleteFinancial(int64Id)
	if err != nil {
		session.AddFlash("Gagal menghapus data keuangan", "error")
		session.Save(request, httpWriter)
	} else {
		session.AddFlash("Berhasil menghapus data keuangan", "success")
		session.Save(request, httpWriter)
	}

	http.Redirect(httpWriter, request, "/home", http.StatusSeeOther)
}

func (controller *FinancialController) DownloadFinancialRecord(httpWriter http.ResponseWriter, request *http.Request) {
	templateLayout := "views/financial/download.html"
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

	selectedMonth := request.URL.Query().Get("selected_month")
	data["selectedMonth"] = selectedMonth

	pemasukanOnly := request.URL.Query().Get("pemasukanOnly") == "true"
	data["pemasukanOnly"] = pemasukanOnly
	pengeluaranOnly := request.URL.Query().Get("pengeluaranOnly") == "true"
	data["pengeluaranOnly"] = pengeluaranOnly

	sessionUserId := session.Values["id"].(int)
	financialModel := models.NewFinancialModel(controller.db)
	total_pemasukan, total_pengeluaran, err := financialModel.GetFinancialTotalNominal(sessionUserId, selectedMonth, pemasukanOnly, pengeluaranOnly)
	if err != nil {
		log.Println("Failed to get financial total nominal data", err)
		data["error"] = "Gagal mendapatkan total data keuangan"
	} else {
		data["total_pemasukan"] = total_pemasukan
		data["total_pengeluaran"] = total_pengeluaran
	}

	financials, err := financialModel.FindAllFinancial(sessionUserId, selectedMonth, pemasukanOnly, pengeluaranOnly)
	if err != nil {
		log.Println("Failed to get all financials data", err)
		data["error"] = "Gagal mendapatkan semua data keuangan"
	} else {
		data["financials"] = financials
	}

	funcMap := template.FuncMap{
		"formatIDR": formatIDR,
		"indexNo": func(a, b int) int { return a + b },
	}
	tmpl, _ := template.New(filepath.Base(templateLayout)).Funcs(funcMap).ParseFiles(templateLayout)
	_ = tmpl.Execute(httpWriter, data)
}