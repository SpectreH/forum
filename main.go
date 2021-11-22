package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	sqlitecommands "forum/sql"
)

type Image struct {
	Name      string
	Type      string
	Container string
}

type Post struct {
	PostId     int
	Author     string
	Title      string
	Body       []string
	Created    string
	Likes      int
	DisLikes   int
	Comments   int
	Categories []string
	Image      Image
}

type Comment struct {
	Author  string
	Created string
	Body    []string
}

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
	Message    string
	AlertType  string
	Posts      []Post
	Categories []string
	LoggedIn   bool
}

type PostPage struct {
	Message   string
	AlertType string
	Post      Post
	Comments  []Comment
	LoggedIn  bool
	Like      bool
	Dislike   bool
}

var MAINPAGEDATA MainPage
var REGDATA Registration
var LOGINDATA Login
var POSTID int

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
		Posts:     nil,
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
	http.HandleFunc("/new", LoadNewPostPage)
	http.HandleFunc("/favicon.ico", Void)

	fmt.Println("Server is listening on port 8000...")
	if http.ListenAndServe(":8000", nil) != nil {
		log.Fatalf("%v - Internal Server Error", http.StatusInternalServerError)
	}
}

func Void(w http.ResponseWriter, r *http.Request) {
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
	if r.URL.Path != "/" {
		if RedirectToPostPage(r.URL.Path) {
			LoadPostPage(w, r)
		} else {
			fmt.Fprint(w, "Error 404")
		}
		return
	}

	templ, _ := template.ParseFiles("templates/main.html")

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

	db, err := sql.Open("sqlite3", "./db/forum.db")
	CheckErr(err)

	MAINPAGEDATA.Posts = CollectAllPostsData(db)
	MAINPAGEDATA.Categories = sqlitecommands.GetAllCategoriesFromTable(db)

	CheckForCookies(r, w)
	if err := templ.Execute(w, MAINPAGEDATA); err != nil {
		panic(err)
	}

	db.Close()

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

		uid, loginExists := sqlitecommands.CheckDataExistence(db, login, "username")
		if !loginExists {
			uid, loginExists = sqlitecommands.CheckDataExistence(db, login, "email")
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
		_, freeUserName := sqlitecommands.CheckDataExistence(db, REGDATA.Username, "username")
		_, freeEmail := sqlitecommands.CheckDataExistence(db, REGDATA.Email, "email")

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

	db, err := sql.Open("sqlite3", "./db/forum.db")
	CheckErr(err)

	if r.Method == "POST" {
		authorId := sqlitecommands.GetUserIdByCookies(db, r, w)
		date := time.Now().Format("2006-01-02 15:04:05")

		if r.FormValue("comment") != "" {
			comment := base64.StdEncoding.EncodeToString([]byte(r.FormValue("comment")))
			sqlitecommands.UpdatePostsCommentsTable(db, POSTID, authorId, date, comment)
		} else {
			var b []byte
			b, _ = ioutil.ReadAll(r.Body)

			if string(b) == "1" { // Like
				sqlitecommands.UpdateRatingsTable(db, POSTID, authorId, "add", "posts_likes")
			} else if string(b) == "2" { // Remove Like
				sqlitecommands.UpdateRatingsTable(db, POSTID, authorId, "remove", "posts_likes")
			} else if string(b) == "-1" { // DisLike
				sqlitecommands.UpdateRatingsTable(db, POSTID, authorId, "add", "posts_dislikes")
			} else if string(b) == "-2" { // Remove DisLike
				sqlitecommands.UpdateRatingsTable(db, POSTID, authorId, "remove", "posts_dislikes")
			}
			db.Close()
			return
		}
	}

	var postPageData PostPage

	postPageData.Post = CollectPostData(db)
	postPageData.LoggedIn = CheckForCookies(r, w)
	if postPageData.LoggedIn == true {
		postPageData.Like = sqlitecommands.GetUserScoreOnPost(db, POSTID, sqlitecommands.GetUserIdByCookies(db, r, w), "posts_likes")
		postPageData.Dislike = sqlitecommands.GetUserScoreOnPost(db, POSTID, sqlitecommands.GetUserIdByCookies(db, r, w), "posts_dislikes")
	} else {
		postPageData.Like = false
		postPageData.Dislike = false
	}

	postPageData.Comments = CollectAllPostComments(db, POSTID)
	db.Close()

	if err := templ.Execute(w, postPageData); err != nil {
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
		authorId, _ := sqlitecommands.CheckDataExistence(db, c.Value, "session_token")

		postTitle := r.FormValue("title")
		postCategories := strings.Split(r.FormValue("categories"), ",")
		postContent := base64.StdEncoding.EncodeToString([]byte(r.FormValue("new-content")))
		date := time.Now().Format("2006-01-02 15:04:05")

		_, data, _ := r.FormFile("myImage")
		postImageData := strings.Split(data.Filename, ".")
		imageContainer := CreateImageContainer(data)

		sqlitecommands.UpdatePostsTable(db, authorId, postTitle, postContent, date, postImageData[0], imageContainer, postImageData[1], postCategories)

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

func RedirectToPostPage(URL string) bool {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	CheckErr(err)

	first, last := sqlitecommands.GetPostsIdGap(db)
	if first == -1 && last == -1 {
		return false
	}

	db.Close()
	postId := strings.Trim(URL, "/")
	number, err := strconv.Atoi(postId)
	if err == nil {
		if number <= last && first <= number {
			POSTID = number
			return true
		}
	}
	return false
}

func CreateImageContainer(file *multipart.FileHeader) string {
	imageByteContainer := make([]byte, (1024 * 1024 * 2))
	fileContent, err := file.Open()

	imageByteContainer, err = ioutil.ReadAll(fileContent)
	if err != nil {
		panic(err)
	}

	fileContent.Close()

	return base64.StdEncoding.EncodeToString(imageByteContainer)
}

func CheckForCookies(r *http.Request, w http.ResponseWriter) bool {
	c, err := r.Cookie("session_token")
	if err == nil {
		db, err := sql.Open("sqlite3", "./db/forum.db")
		CheckErr(err)

		_, checkResult := sqlitecommands.CheckDataExistence(db, c.Value, "session_token")
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
		Expires: time.Now().Add(1200 * time.Second),
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

func DivideBodyIntoParagraphs(body string) []string {
	var result []string
	var paragraph []byte

	base64Body, err := base64.StdEncoding.DecodeString(body)
	CheckErr(err)

	for i := 0; i < len(base64Body); i++ {
		paragraph = append(paragraph, base64Body[i])

		if base64Body[i] == 13 {
			result = append(result, string(paragraph))
			i = i + 2
			paragraph = make([]byte, 0)
		}

		if len(paragraph) != 0 && i == len(base64Body)-1 {
			result = append(result, string(paragraph))
		}
	}

	return result
}

func CollectPostData(db *sql.DB) Post {
	var post Post
	var date time.Time
	var body string
	var authorId int

	_ = db.QueryRow("SELECT * FROM posts WHERE id = ?", POSTID).Scan(&post.PostId, &authorId, &post.Title, &body, &date, &post.Likes, &post.DisLikes, &post.Comments)

	post.Created = date.Format("January 02, 2006 at 15:04")
	post.Author = sqlitecommands.GetUserNameFromTable(db, authorId)
	post.Categories = sqlitecommands.GetPostCategoriesFromTable(db, post.PostId)
	post.Body = DivideBodyIntoParagraphs(body)

	return post
}

func CollectAllPostsData(db *sql.DB) []Post {
	var result []Post

	rows, err := db.Query("SELECT * FROM posts")
	CheckErr(err)

	for rows.Next() {
		var post Post
		var date time.Time
		var body string
		var authorId int

		err := rows.Scan(&post.PostId, &authorId, &post.Title, &body, &date, &post.Likes, &post.DisLikes, &post.Comments)
		if err != nil {
			log.Fatal(err)
		}

		post.Created = date.Format("January 02, 2006 at 15:04")
		post.Author = sqlitecommands.GetUserNameFromTable(db, authorId)
		post.Categories = sqlitecommands.GetPostCategoriesFromTable(db, post.PostId)
		post.Image.Name, post.Image.Container, post.Image.Type = sqlitecommands.GetImageDataFromTable(db, post.PostId)
		post.Body = DivideBodyIntoParagraphs(body)

		result = append(result, post)
	}

	return result
}

func CollectAllPostComments(db *sql.DB, postId int) []Comment {
	var result []Comment

	rows, err := db.Query("SELECT * FROM posts_comments WHERE post_id = ?", postId)
	CheckErr(err)

	for rows.Next() {
		var comment Comment
		var body string
		var date time.Time
		var id, postId, authorId int

		err := rows.Scan(&id, &postId, &authorId, &date, &body)
		if err != nil {
			log.Fatal(err)
		}

		comment.Author = sqlitecommands.GetUserNameFromTable(db, authorId)
		comment.Created = date.Format("January 02, 2006 at 15:04")
		comment.Body = DivideBodyIntoParagraphs(body)

		result = append(result, comment)
	}

	return result
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
