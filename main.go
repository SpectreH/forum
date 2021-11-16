package main

import (
	"database/sql"
	"fmt"
	sqlitecommands "forum/sql"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
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
var REGDATA Registration
var LOGINDATA Login

func main() {
	REGDATA := Registration{
		NameErr:  false,
		EmailErr: false,
		Username: "nil",
		Email:    "nil",
	}

	LOGINDATA := Login{
		LoginErr: false,
		PassErr:  false,
		Login:    "",
	}

	MAINPAGEDATA = MainPage{
		Message:   "",
		AlertType: "",
		LoggedIn:  false,
	}

	_ = LOGINDATA
	_ = REGDATA

	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	js := http.FileServer(http.Dir("js"))
	http.Handle("/js/", http.StripPrefix("/js/", js))

	http.HandleFunc("/", LoadMainPage)
	http.HandleFunc("/logout", LogOut)
	http.HandleFunc("/login", LoadLoginPage)
	http.HandleFunc("/registration", LoadRegistrationPage)
	http.HandleFunc("/exit", ShutdownServer)
	http.HandleFunc("/1", LoadPostPage)
	http.HandleFunc("/new", LoadNewPostPage)

	fmt.Println("Server is listening on port 8000...")
	if http.ListenAndServe(":8000", nil) != nil {
		log.Fatalf("%v - Internal Server Error", http.StatusInternalServerError)
	}
}

func LogOut(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, "/", 302)
}

func LoadMainPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("templates/main.html")

	CheckForCookies(r, w)

	REGDATA = Registration{
		NameErr:  false,
		EmailErr: false,
		Username: "nil",
		Email:    "nil",
	}

	LOGINDATA = Login{
		LoginErr: false,
		PassErr:  false,
		Login:    "",
	}

	if err := templ.Execute(w, MAINPAGEDATA); err != nil {
		panic(err)
	}

	MAINPAGEDATA = MainPage{
		Message:   "",
		AlertType: "",
	}
}

func LoadLoginPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("templates/login.html")

	if r.Method == "GET" {
		if CheckForCookies(r, w) {
			RedirectToMainPage(r, w, "You are already logged in!", "AlreadyLoged")
		}
	}

	if r.Method == "POST" {
		login := r.FormValue("login")
		password := []byte(r.FormValue("password"))

		db, err := sql.Open("sqlite3", "./db/forum.db")
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
				LOGINDATA.LoginErr = false
				LOGINDATA.PassErr = false

				sqlitecommands.UpdateSessionToken(db, CreateSessionToken(w), uid)

				RedirectToMainPage(r, w, "Successfully logged in!", "Login")
				db.Close()
			} else {
				LOGINDATA.Login = login
				LOGINDATA.LoginErr = false
				LOGINDATA.PassErr = true
			}
		} else {
			LOGINDATA.Login = ""
			LOGINDATA.LoginErr = true
			LOGINDATA.PassErr = false
		}
	}

	if err := templ.Execute(w, LOGINDATA); err != nil {
		panic(err)
	}
}

func LoadRegistrationPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("templates/registration.html")

	if r.Method == "GET" {
		if CheckForCookies(r, w) {
			RedirectToMainPage(r, w, "You are already registered and logged in!", "AlreadyRegistered")
		}
	}

	REGDATA.NameErr = false
	REGDATA.EmailErr = false

	if r.Method == "POST" {
		REGDATA.Username = r.FormValue("username")
		REGDATA.Email = r.FormValue("email")
		password := GetHash([]byte(r.FormValue("password")))
		date := time.Now().Format("2006-01-02 15:04:05")
		role := 1
		ip := "0"

		db, err := sql.Open("sqlite3", "./db/forum.db")
		CheckErr(err)

		// Checks if REGDATA is already taken
		_, freeUserName := DataExists(db, REGDATA.Username, "username")
		_, freeEmail := DataExists(db, REGDATA.Email, "email")

		if freeUserName || freeEmail {
			if freeUserName {
				REGDATA.NameErr = true
			}
			if freeEmail {
				REGDATA.EmailErr = true
			}
		} else {
			sqlitecommands.UpdateUsersTable(db, CreateSessionToken(w), REGDATA.Username, REGDATA.Email, password, date, role, ip)

			RedirectToMainPage(r, w, "Account successfully created!", "Register")
			db.Close()
		}
	}

	if err := templ.Execute(w, REGDATA); err != nil {
		panic(err)
	}
}

func LoadPostPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("templates/post.html")

	fmt.Println(r.Method)

	CheckForCookies(r, w)

	if err := templ.Execute(w, MAINPAGEDATA); err != nil {
		panic(err)
	}
}

func LoadNewPostPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("templates/new.html")

	if !CheckForCookies(r, w) {
		RedirectToMainPage(r, w, "You are not logged in!", "NotLoggedIn")
		return
	}

	if r.Method == "POST" {
		db, err := sql.Open("sqlite3", "./db/forum.db")
		CheckErr(err)

		c, _ := r.Cookie("session_token")
		authorId, _ := DataExists(db, c.Value, "session_token")

		postTitle := r.FormValue("title")
		postCategories := strings.Split(r.FormValue("categories"), ",")
		postContent := r.FormValue("new-content")
		date := time.Now().Format("2006-01-02 15:04:05")

		_, data, _ := r.FormFile("myImage")
		postImageData := strings.Split(data.Filename, ".")
		imageByteContainer := CreateImageContainer(data)

		sqlitecommands.UpdatePostsTable(db, authorId, postTitle, postContent, date, postImageData[0], imageByteContainer, postImageData[1], postCategories)

		db.Close()
	}

	CheckForCookies(r, w)

	if err := templ.Execute(w, MAINPAGEDATA); err != nil {
		panic(err)
	}
}

func RedirectToMainPage(r *http.Request, w http.ResponseWriter, message string, alertType string) {
	MAINPAGEDATA = MainPage{
		Message:   message,
		AlertType: alertType,
	}

	http.Redirect(w, r, "/", 302)
}

func CreateImageContainer(file *multipart.FileHeader) []byte {
	imageByteContainer := make([]byte, (1024 * 1024 * 2))
	fileContent, err := file.Open()

	imageByteContainer, err = ioutil.ReadAll(fileContent)
	if err != nil {
		panic(err)
	}

	fileContent.Close()

	return imageByteContainer
}

func CheckForCookies(r *http.Request, w http.ResponseWriter) bool {
	c, err := r.Cookie("session_token")

	if err == nil {
		db, err := sql.Open("sqlite3", "./db/forum.db")
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

func DataExists(db *sql.DB, REGDATA string, dataType string) (int, bool) {
	var uid int
	sqlStmt := "SELECT " + dataType + ", uid FROM users WHERE " + dataType + " = ?"
	err := db.QueryRow(sqlStmt, REGDATA).Scan(&REGDATA, &uid)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return -1, false
	}
	return uid, true
}

func ShutdownServer(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
