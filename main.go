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
	"golang.org/x/crypto/bcrypt"
)

type Registration struct {
	Success  bool
	NameErr  bool
	EmailErr bool
	Username string
	Email    string
}

func main() {
	data := Registration{
		Success:  false,
		NameErr:  false,
		EmailErr: false,
		Username: "nil",
		Email:    "nil",
	}

	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	http.HandleFunc("/", LoadMainPage(data))
	http.HandleFunc("/login", LoadLoginPage(data))
	http.HandleFunc("/registration", LoadRegistrationPage(data))
	http.HandleFunc("/exit", ShutdownServer)

	fmt.Println("Server is listening on port 8000...")
	if http.ListenAndServe(":8000", nil) != nil {
		log.Fatalf("%v - Internal Server Error", http.StatusInternalServerError)
	}
}

func LoadMainPage(data Registration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, _ := template.ParseFiles("templates/index.html")
		ExecuteTempl(templ, w, data)
	}
}

func LoadLoginPage(data Registration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, _ := template.ParseFiles("templates/login.html")
		ExecuteTempl(templ, w, data)
	}
}

func LoadRegistrationPage(data Registration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ, _ := template.ParseFiles("templates/registration.html")

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
			freeUserName := DataExists(db, data.Username, "username")
			freeEmail := DataExists(db, data.Email, "email")

			if freeUserName || freeEmail {
				if freeUserName {
					data.NameErr = true
				}
				if freeEmail {
					data.EmailErr = true
				}
				data.Success = false
			} else {
				stmt, err := db.Prepare("INSERT INTO users(username, email, password, date, role, ip) values(?,?,?,?,?,?)")
				CheckErr(err)

				_, err = stmt.Exec(data.Username, data.Email, password, date, role, ip)
				CheckErr(err)

				data.Success = true
				http.Redirect(w, r, "/", 302)
				db.Close()
			}
		}

		ExecuteTempl(templ, w, data)
	}
}

func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func DataExists(db *sql.DB, data string, dataType string) bool {
	sqlStmt := "SELECT " + dataType + " FROM users WHERE " + dataType + " = ?"
	err := db.QueryRow(sqlStmt, data).Scan(&data)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return false
	}
	return true
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
