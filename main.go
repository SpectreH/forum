package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Registration struct {
	NameErr  bool
	EmailErr bool
	Username string
	Email    string
}

type Login struct {
	LoginErr bool
	PassErr  bool
	Login    string
}

type MainPage struct {
	Message   string
	AlertType string
	LoggedIn  bool
}

var MAINPAGEDATA MainPage

func main() {
	data := Registration{
		NameErr:  false,
		EmailErr: false,
		Username: "nil",
		Email:    "nil",
	}

	loginData := Login{
		LoginErr: false,
		PassErr:  false,
		Login:    "",
	}

	MAINPAGEDATA = MainPage{
		Message:   "",
		AlertType: "",
		LoggedIn:  false,
	}

	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	js := http.FileServer(http.Dir("js"))
	http.Handle("/js/", http.StripPrefix("/js/", js))

	http.HandleFunc("/", LoadMainPage())
	http.HandleFunc("/logout", LogOut())
	http.HandleFunc("/login", LoadLoginPage(&loginData))
	http.HandleFunc("/registration", LoadRegistrationPage(&data))
	http.HandleFunc("/exit", ShutdownServer)

	fmt.Println("Server is listening on port 8000...")
	if http.ListenAndServe(":8000", nil) != nil {
		log.Fatalf("%v - Internal Server Error", http.StatusInternalServerError)
	}
}

func LogOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if CheckForCookies(r, w) {
			c := http.Cookie{
				Name:   "session_token",
				MaxAge: -1}
			http.SetCookie(w, &c)

			MAINPAGEDATA = MainPage{
				Message:   "You have successfully logged out!",
				AlertType: "Logout",
			}
		}

		templ, _ := template.ParseFiles("templates/main.html")
		if err := templ.Execute(w, MAINPAGEDATA); err != nil {
			panic(err)
		}

		MAINPAGEDATA = MainPage{
			Message:   "",
			AlertType: "",
		}
	}
}

func LoadMainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, _ := template.ParseFiles("templates/main.html")

		CheckForCookies(r, w)

		if err := templ.Execute(w, MAINPAGEDATA); err != nil {
			panic(err)
		}

		MAINPAGEDATA = MainPage{
			Message:   "",
			AlertType: "",
		}
	}
}

func LoadLoginPage(loginData *Login) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, _ := template.ParseFiles("templates/login.html")

		if r.Method == "GET" {
			if CheckForCookies(r, w) {
				MAINPAGEDATA = MainPage{
					Message:   "You are already logged in!",
					AlertType: "AlreadyLoged",
				}

				http.Redirect(w, r, "/", 302)
			}
		}

		if r.Method == "POST" {
			login := r.FormValue("login")
			password := []byte(r.FormValue("password"))

			db, err := sql.Open("sqlite3", "./db/users.db")
			CheckErr(err)

			uid, loginExists := DataExists(db, login, "username")
			if !loginExists {
				uid, loginExists = DataExists(db, login, "email")
			}

			if loginExists {
				var accountHash string

				sqlStmt := "SELECT password FROM users WHERE username = ? OR email = ?"
				_ = db.QueryRow(sqlStmt, login, login).Scan(&accountHash)

				if bcrypt.CompareHashAndPassword([]byte(accountHash), password) == nil {
					loginData.LoginErr = false
					loginData.PassErr = false

					sessionToken := CreateSessionToken(w)

					stmt, err := db.Prepare("UPDATE users SET session_token = ? WHERE uid = ?")
					CheckErr(err)
					_, err = stmt.Exec(sessionToken, uid)
					CheckErr(err)

					MAINPAGEDATA = MainPage{
						Message:   "Successfully logged in!",
						AlertType: "Login",
					}

					http.Redirect(w, r, "/", 302)
					db.Close()
				} else {
					loginData.Login = login
					loginData.LoginErr = false
					loginData.PassErr = true
				}
			} else {
				loginData.Login = ""
				loginData.LoginErr = true
				loginData.PassErr = false
			}
		}

		if err := templ.Execute(w, loginData); err != nil {
			panic(err)
		}
	}
}

func LoadRegistrationPage(data *Registration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, _ := template.ParseFiles("templates/registration.html")

		if r.Method == "GET" {
			if CheckForCookies(r, w) {
				MAINPAGEDATA = MainPage{
					Message:   "You are already registered and logged in!",
					AlertType: "AlreadyRegistered",
				}

				http.Redirect(w, r, "/", 302)
			}
		}

		data.NameErr = false
		data.EmailErr = false

		if r.Method == "POST" {
			data.Username = r.FormValue("username")
			data.Email = r.FormValue("email")
			password := GetHash([]byte(r.FormValue("password")))
			date := time.Now().Format("2006-01-02 15:04:05")
			role := 1
			ip := 0

			db, err := sql.Open("sqlite3", "./db/users.db")
			CheckErr(err)

			// Checks if data is already taken
			_, freeUserName := DataExists(db, data.Username, "username")
			_, freeEmail := DataExists(db, data.Email, "email")

			if freeUserName || freeEmail {
				if freeUserName {
					data.NameErr = true
				}
				if freeEmail {
					data.EmailErr = true
				}
			} else {
				stmt, err := db.Prepare("INSERT INTO users(username, email, password, date, role, ip) values(?,?,?,?,?,?)")
				CheckErr(err)

				_, err = stmt.Exec(data.Username, data.Email, password, date, role, ip)
				CheckErr(err)

				sessionToken := CreateSessionToken(w)

				stmt, err = db.Prepare("UPDATE users SET session_token = ? WHERE username = ?")
				CheckErr(err)
				_, err = stmt.Exec(sessionToken, data.Username)
				CheckErr(err)

				MAINPAGEDATA = MainPage{
					Message:   "Account successfully created!",
					AlertType: "Register",
				}

				http.Redirect(w, r, "/", 302)
				db.Close()
			}
		}

		ExecuteTempl(templ, w, *data)
	}
}

func CheckForCookies(r *http.Request, w http.ResponseWriter) bool {
	c, err := r.Cookie("session_token")

	if err == nil {
		db, err := sql.Open("sqlite3", "./db/users.db")
		CheckErr(err)

		_, checkResult := DataExists(db, c.Value, "session_token")
		db.Close()

		if checkResult {
			MAINPAGEDATA.LoggedIn = true
			return true
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	MAINPAGEDATA.LoggedIn = false
	return false
}

func CreateSessionToken(w http.ResponseWriter) string {
	sessionToken := uuid.NewV4().String()
	MAINPAGEDATA.LoggedIn = true

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(20 * time.Second),
	})
	return sessionToken
}

func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func DataExists(db *sql.DB, data string, dataType string) (int, bool) {
	var uid int
	sqlStmt := "SELECT " + dataType + ", uid FROM users WHERE " + dataType + " = ?"
	err := db.QueryRow(sqlStmt, data).Scan(&data, &uid)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return -1, false
	}
	return uid, true
}

func ExecuteTempl(templ *template.Template, w http.ResponseWriter, data Registration) {
	if err := templ.Execute(w, data); err != nil {
		panic(err)
	}
}

func ShutdownServer(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
