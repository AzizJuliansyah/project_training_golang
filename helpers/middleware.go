package helpers

import (
	"financial_record/config"
	"financial_record/views"
	"net/http"
)

func GuestOnly(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, config.SESSION_ID)
		if session.Values["loggedIn"] == true {
			if session.Values["isAdmin"] == true {
				http.Redirect(w, r, "/home-admin", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/home", http.StatusSeeOther)
			}
			return
		}

		next.ServeHTTP(w, r)
	}
}

func AuthOnly(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, config.SESSION_ID)
		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, config.SESSION_ID)
		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		isAdmin := session.Values["isAdmin"]
		if isAdmin == false {
			data := map[string]interface{}{
				"isAdmin": isAdmin,
			}
			views.RenderTemplate(w, "views/static/forbidden/forbidden.html", data)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func EmployeeOnly(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := config.Store.Get(r, config.SESSION_ID)
		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		isAdmin := session.Values["isAdmin"]
		if isAdmin == true {
			data := map[string]interface{}{
				"isAdmin": isAdmin,
			}
			views.RenderTemplate(w, "views/static/forbidden/forbidden.html", data)
			return
		}

		next.ServeHTTP(w, r)
	}
}


