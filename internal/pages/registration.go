package pages

import (
	"database/sql"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"html/template"
	"net/http"
	"regexp"
	"time"
)

type Registration struct {
	DB *sql.DB
}

type RegistrationData struct {
	NameErr  string
	EmailErr string
	Username string
	Email    string
}

func (data Registration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("templates/registration.html")

	if r.Method == "GET" {
		if utility.CheckForCookies(data.DB, r, w) {
			utility.RedirectToMainPage(r, w, "You are already registered and logged in!", "AlreadyRegistered")
		}
	}

	registrationData := RegistrationData{
		NameErr:  "",
		EmailErr: "",
		Username: "",
		Email:    "",
	}

	if r.Method == "POST" {
		registrationData.Username = r.FormValue("username")
		registrationData.Email = r.FormValue("email")
		password := utility.GetHash([]byte(r.FormValue("password")))
		date := time.Now().Format("2006-01-02 15:04:05")
		role := 1
		ip := "0"

		// Checks if username or email are already taken
		_, freeUserName := sqlitecommands.CheckDataExistence(data.DB, registrationData.Username, "username")
		_, freeEmail := sqlitecommands.CheckDataExistence(data.DB, registrationData.Email, "email")

		if !ValidateUserNameInput(registrationData.Username) {
			registrationData.NameErr = "Only letters and numbers are allowed"
			registrationData.Username = ""
		}

		if !ValidateEmailInput(registrationData.Email) {
			registrationData.EmailErr = "Email address format should be example@mail.com"
			registrationData.Email = ""
		}

		if freeUserName {
			registrationData.NameErr = "Username is already taken"
			registrationData.Username = ""
		}

		if freeEmail {
			registrationData.EmailErr = "Email is already registered"
			registrationData.Email = ""
		}

		if registrationData.EmailErr == "" && registrationData.NameErr == "" {
			sqlitecommands.UpdateUsersTable(data.DB, utility.CreateSessionToken(w), registrationData.Username, registrationData.Email, password, date, role, ip)
			utility.RedirectToMainPage(r, w, "Account successfully created!", "Register")
		}
	}

	if err := templ.Execute(w, registrationData); err != nil {
		panic(err)
	}
}

func ValidateUserNameInput(value string) bool {
	usernameFormat := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
	return usernameFormat.MatchString(value)
}

func ValidateEmailInput(value string) bool {
	emailFormat := regexp.MustCompile(`.+@+.+\..+`)
	return emailFormat.MatchString(value)
}
