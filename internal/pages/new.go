package pages

import (
	"encoding/base64"
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"net/http"
	"strings"
	"time"
)

type New struct {
}

func (data New) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !utility.CheckForCookies(r, w) {
		utility.RedirectToMainPage(r, w, "You are not logged in!", "Fail_NotLoggedIn")
		return
	}

	if r.Method == "POST" {
		c, _ := r.Cookie("session_token")
		authorId, _ := sqlitecommands.CheckDataExistence(c.Value, "session_token")
		postTitle := r.FormValue("title")
		postCategories := strings.Split(r.FormValue("categories"), ",")
		postContent := base64.StdEncoding.EncodeToString([]byte(r.FormValue("new-content")))
		date := time.Now().Format("2006-01-02 15:04:05")
		imgPath := utility.SavePostImage(r)

		sqlitecommands.UpdatePostsTable(authorId, postTitle, postContent, date, imgPath, postCategories)
		utility.RedirectToMainPage(r, w, "You have successfully added a new post!", "Success")
		return
	}

	if err := env.TEMPLATES["new.html"].Execute(w, env.MAINPAGEDATA); err != nil {
		panic(err)
	}
}
