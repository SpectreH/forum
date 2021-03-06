package pages

import (
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"net/http"
	"regexp"
	"time"
)

type Registration struct {
}

type RegistrationData struct {
	NameErr  string
	EmailErr string
	Username string
	Email    string
}

func (data Registration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if utility.CheckForCookies(r, w) {
			utility.RedirectToMainPage(r, w, "You are already registered and logged in!", "Fail_AlreadyRegistered")
		}
	}

	registrationData := RegistrationData{}

	if r.Method == "POST" {
		registrationData.Username = r.FormValue("username")
		registrationData.Email = r.FormValue("email")
		password := utility.GetHash([]byte(r.FormValue("password")))
		date := time.Now().Format("2006-01-02 15:04:05")
		role := 1
		ip := "0"

		// Checks if username or email are already taken
		_, freeUserName := sqlitecommands.CheckDataExistence(registrationData.Username, "username")
		_, freeEmail := sqlitecommands.CheckDataExistence(registrationData.Email, "email")

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
			sqlitecommands.UpdateUsersTable(utility.CreateSessionToken(w), registrationData.Username, registrationData.Email, password, date, role, ip)
			utility.RedirectToMainPage(r, w, "Account successfully created!", "Success")
		}
	}

	if err := env.TEMPLATES["registration.html"].Execute(w, registrationData); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
