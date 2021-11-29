package pages

import (
	"database/sql"
	"encoding/base64"
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type New struct {
	DB *sql.DB
}

func (data New) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("templates/new.html")

	if !utility.CheckForCookies(data.DB, r, w) {
		utility.RedirectToMainPage(r, w, "You are not logged in!", "NotLoggedIn")
		return
	}

	if r.Method == "POST" {
		c, _ := r.Cookie("session_token")
		authorId, _ := sqlitecommands.CheckDataExistence(data.DB, c.Value, "session_token")

		postTitle := r.FormValue("title")
		postCategories := strings.Split(r.FormValue("categories"), ",")
		postContent := base64.StdEncoding.EncodeToString([]byte(r.FormValue("new-content")))
		date := time.Now().Format("2006-01-02 15:04:05")

		_, imageData, _ := r.FormFile("myImage")
		postImageData := strings.Split(imageData.Filename, ".")
		imageContainer := utility.CreateImageContainer(imageData)

		sqlitecommands.UpdatePostsTable(data.DB, authorId, postTitle, postContent, date, postImageData[0], imageContainer, postImageData[1], postCategories)

		utility.RedirectToMainPage(r, w, "You have successfully added a new post!", "newPost")
		return
	}

	if err := templ.Execute(w, env.MAINPAGEDATA); err != nil {
		panic(err)
	}
}
