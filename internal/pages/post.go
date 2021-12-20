package pages

import (
	"encoding/base64"
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func LoadPostPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if !utility.CheckForCookies(r, w) {
			return
		}

		authorId := sqlitecommands.GetUserIdByCookies(r, w)
		date := time.Now().Format("2006-01-02 15:04:05")

		if r.FormValue("comment") != "" {
			if r.FormValue("comment") == "" {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			comment := base64.StdEncoding.EncodeToString([]byte(r.FormValue("comment")))
			sqlitecommands.UpdatePostsCommentsTable(env.POSTID, authorId, date, comment)

			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		} else {
			var b []byte
			b, _ = ioutil.ReadAll(r.Body)

			rateDataSep := strings.Split(string(b), ";")
			id, _ := strconv.Atoi(rateDataSep[2])

			if rateDataSep[0] == "post" {
				if string(rateDataSep[1]) == "1" { // Like
					sqlitecommands.UpdatePostRatingTable(id, authorId, "add", "posts_likes")
					sqlitecommands.UpdatePostRatingTable(id, authorId, "remove", "posts_dislikes")
				} else if string(rateDataSep[1]) == "2" { // Remove Like
					sqlitecommands.UpdatePostRatingTable(id, authorId, "remove", "posts_likes")
				} else if string(rateDataSep[1]) == "-1" { // DisLike
					sqlitecommands.UpdatePostRatingTable(id, authorId, "add", "posts_dislikes")
					sqlitecommands.UpdatePostRatingTable(id, authorId, "remove", "posts_likes")
				} else if string(rateDataSep[1]) == "-2" { // Remove DisLike
					sqlitecommands.UpdatePostRatingTable(id, authorId, "remove", "posts_dislikes")
				}
			} else {
				if string(rateDataSep[1]) == "1" { // Like
					sqlitecommands.UpdateCommentRatingTable(id, authorId, "add", "comment_likes")
					sqlitecommands.UpdateCommentRatingTable(id, authorId, "remove", "comment_dislikes")
				} else if string(rateDataSep[1]) == "2" { // Remove Like
					sqlitecommands.UpdateCommentRatingTable(id, authorId, "remove", "comment_likes")
				} else if string(rateDataSep[1]) == "-1" { // DisLike
					sqlitecommands.UpdateCommentRatingTable(id, authorId, "add", "comment_dislikes")
					sqlitecommands.UpdateCommentRatingTable(id, authorId, "remove", "comment_likes")
				} else if string(rateDataSep[1]) == "-2" { // Remove DisLike
					sqlitecommands.UpdateCommentRatingTable(id, authorId, "remove", "comment_dislikes")
				}
			}
			return
		}
	}

	var postPageData env.PostPage
	postPageData.Post = utility.CollectPostData()
	postPageData.LoggedIn = utility.CheckForCookies(r, w)

	if postPageData.LoggedIn {
		postPageData.Post.Liked = sqlitecommands.GetUserScoreOnPost(env.POSTID, sqlitecommands.GetUserIdByCookies(r, w), "posts_likes")
		postPageData.Post.Disliked = sqlitecommands.GetUserScoreOnPost(env.POSTID, sqlitecommands.GetUserIdByCookies(r, w), "posts_dislikes")
	} else {
		postPageData.Post.Liked = false
		postPageData.Post.Disliked = false
	}

	postPageData.Comments = utility.CollectAllPostComments(env.POSTID, w, r)

	if err := env.TEMPLATES["post.html"].Execute(w, postPageData); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
