package pages

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"io"
	"log"
	"net/http"
	"os"
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
		imgPath := SavePostImage(r)

		sqlitecommands.UpdatePostsTable(authorId, postTitle, postContent, date, imgPath, postCategories)
		utility.RedirectToMainPage(r, w, "You have successfully added a new post!", "Success")
		return
	}

	if err := env.TEMPLATES["new.html"].Execute(w, env.MAINPAGEDATA); err != nil {
		panic(err)
	}
}

func SavePostImage(r *http.Request) string {
	var path string

	in, header, err := r.FormFile("myImage")
	imageData := strings.Split(header.Filename, ".")
	if err != nil {
		log.Println(err)
	}
	defer in.Close()

	randBytes := make([]byte, 16)
	rand.Read(randBytes)

	path = "images/" + hex.EncodeToString(randBytes) + "." + imageData[1]

	out, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()
	io.Copy(out, in)

	return path
}
