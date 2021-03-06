package pages

import (
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
}

type LoginData struct {
	LoginErr string
	PassErr  string
	Login    string
}

func (data Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if utility.CheckForCookies(r, w) {
			utility.RedirectToMainPage(r, w, "You are already logged in!", "Fail_AlreadyLogged")
		}
	}

	loginData := LoginData{}

	if r.Method == "POST" {
		login := r.FormValue("login")
		password := []byte(r.FormValue("password"))

		uid, dataExists := sqlitecommands.CheckDataExistence(login, "username")

		// Checks email if username doesn't exist
		if !dataExists {
			uid, dataExists = sqlitecommands.CheckDataExistence(login, "email")
		}

		if dataExists {
			var accountHash string
			sqlStmt := "SELECT password FROM users WHERE username = ? OR email = ?"
			_ = env.DB.QueryRow(sqlStmt, login, login).Scan(&accountHash)

			if bcrypt.CompareHashAndPassword([]byte(accountHash), password) == nil {
				sqlitecommands.UpdateSessionToken(utility.CreateSessionToken(w), uid)
				utility.RedirectToMainPage(r, w, "Successfully logged in!", "Success")
			} else {
				loginData.Login = login
				loginData.PassErr = "Password does not match"
			}

		} else {
			loginData.LoginErr = "Account does not exist"
		}
	}

	if err := env.TEMPLATES["login.html"].Execute(w, loginData); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
