package config

import "github.com/gorilla/sessions"

const SESSION_ID = "expense-records-app"

var Store *sessions.CookieStore

func init() {
	Store = sessions.NewCookieStore([]byte(SESSION_ID))
	Store.Options = &sessions.Options{
		Path: "/",
		MaxAge: 3600 * 8,
		HttpOnly: true,
		Secure: true,
	}
}