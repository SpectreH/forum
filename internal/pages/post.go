package pages

import (
	"database/sql"
	"encoding/base64"
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func LoadPostPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("templates/post.html")

	db, err := sql.Open("sqlite3", "./db/forum.db")
	utility.CheckErr(err)

	if r.Method == "POST" {
		if !utility.CheckForCookies(db, r, w) {
			return
		}

		authorId := sqlitecommands.GetUserIdByCookies(db, r, w)
		date := time.Now().Format("2006-01-02 15:04:05")

		if r.FormValue("comment") != "" {
			comment := base64.StdEncoding.EncodeToString([]byte(r.FormValue("comment")))
			sqlitecommands.UpdatePostsCommentsTable(db, env.POSTID, authorId, date, comment)
		} else {
			var b []byte
			b, _ = ioutil.ReadAll(r.Body)

			rateDataSep := strings.Split(string(b), ";")
			id, _ := strconv.Atoi(rateDataSep[2])

			if rateDataSep[0] == "post" {
				if string(rateDataSep[1]) == "1" { // Like
					sqlitecommands.UpdatePostRatingsTable(db, id, authorId, "add", "posts_likes")
				} else if string(rateDataSep[1]) == "2" { // Remove Like
					sqlitecommands.UpdatePostRatingsTable(db, id, authorId, "remove", "posts_likes")
				} else if string(rateDataSep[1]) == "-1" { // DisLike
					sqlitecommands.UpdatePostRatingsTable(db, id, authorId, "add", "posts_dislikes")
				} else if string(rateDataSep[1]) == "-2" { // Remove DisLike
					sqlitecommands.UpdatePostRatingsTable(db, id, authorId, "remove", "posts_dislikes")
				}
			} else {
				if string(rateDataSep[1]) == "1" { // Like
					sqlitecommands.UpdateCommentRatingsTable(db, id, authorId, "add", "comment_likes")
				} else if string(rateDataSep[1]) == "2" { // Remove Like
					sqlitecommands.UpdateCommentRatingsTable(db, id, authorId, "remove", "comment_likes")
				} else if string(rateDataSep[1]) == "-1" { // DisLike
					sqlitecommands.UpdateCommentRatingsTable(db, id, authorId, "add", "comment_dislikes")
				} else if string(rateDataSep[1]) == "-2" { // Remove DisLike
					sqlitecommands.UpdateCommentRatingsTable(db, id, authorId, "remove", "comment_dislikes")
				}
			}

			db.Close()
			return
		}
	}

	var postPageData env.PostPage

	postPageData.Post = utility.CollectPostData(db)
	postPageData.LoggedIn = utility.CheckForCookies(db, r, w)
	if postPageData.LoggedIn == true {
		postPageData.Post.Liked = sqlitecommands.GetUserScoreOnPost(db, env.POSTID, sqlitecommands.GetUserIdByCookies(db, r, w), "posts_likes")
		postPageData.Post.Disliked = sqlitecommands.GetUserScoreOnPost(db, env.POSTID, sqlitecommands.GetUserIdByCookies(db, r, w), "posts_dislikes")
	} else {
		postPageData.Post.Liked = false
		postPageData.Post.Disliked = false
	}

	postPageData.Comments = utility.CollectAllPostComments(db, env.POSTID, w, r)
	db.Close()

	if err := templ.Execute(w, postPageData); err != nil {
		panic(err)
	}
}
