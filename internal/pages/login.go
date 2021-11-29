package pages

import (
	"database/sql"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	DB *sql.DB
}

type LoginData struct {
	LoginErr string
	PassErr  string
	Login    string
}

func (data Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("templates/login.html")

	if r.Method == "GET" {
		if utility.CheckForCookies(data.DB, r, w) {
			utility.RedirectToMainPage(r, w, "You are already logged in!", "AlreadyLoged")
		}
	}

	loginData := LoginData{}

	if r.Method == "POST" {
		login := r.FormValue("login")
		password := []byte(r.FormValue("password"))

		uid, loginExists := sqlitecommands.CheckDataExistence(data.DB, login, "username")
		if !loginExists {
			uid, loginExists = sqlitecommands.CheckDataExistence(data.DB, login, "email")
		}

		if loginExists {
			var accountHash string
			sqlStmt := "SELECT password FROM users WHERE username = ? OR email = ?"
			_ = data.DB.QueryRow(sqlStmt, login, login).Scan(&accountHash)

			if bcrypt.CompareHashAndPassword([]byte(accountHash), password) == nil {
				loginData.LoginErr = ""
				loginData.PassErr = ""

				sqlitecommands.UpdateSessionToken(data.DB, utility.CreateSessionToken(w), uid)

				utility.RedirectToMainPage(r, w, "Successfully logged in!", "Login")
			} else {
				loginData.Login = login
				loginData.LoginErr = ""
				loginData.PassErr = "Password does not match"
			}
		} else {
			loginData.Login = ""
			loginData.LoginErr = "Account does not exist"
			loginData.PassErr = ""
		}
	}

	if err := templ.Execute(w, loginData); err != nil {
		panic(err)
	}
}
