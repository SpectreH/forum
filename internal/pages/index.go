package pages

import (
	"database/sql"
	"fmt"
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"html/template"
	"net/http"
)

type Main struct {
	DB *sql.DB
}

func (data Main) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if utility.RedirectToPostPage(r.URL.Path) {
			LoadPostPage(w, r)
		} else {
			fmt.Fprint(w, "Error 404")
		}
		return
	}

	templ, _ := template.ParseFiles("templates/main.html")

	env.MAINPAGEDATA.Posts = utility.CollectAllPostsData(data.DB)
	env.MAINPAGEDATA.Categories = sqlitecommands.GetAllCategoriesFromTable(data.DB)

	env.MAINPAGEDATA.LoggedIn = utility.CheckForCookies(data.DB, r, w)
	env.MAINPAGEDATA.Username = "-1"

	var userId int
	if env.MAINPAGEDATA.LoggedIn {
		userId = sqlitecommands.GetUserIdByCookies(data.DB, r, w)
		env.MAINPAGEDATA.Username = sqlitecommands.GetUserNameFromTable(data.DB, userId)
	}
	for i := 0; i < len(env.MAINPAGEDATA.Posts); i++ {
		if env.MAINPAGEDATA.LoggedIn == true {
			env.MAINPAGEDATA.Posts[i].Liked = sqlitecommands.GetUserScoreOnPost(data.DB, env.MAINPAGEDATA.Posts[i].PostId, userId, "posts_likes")
			env.MAINPAGEDATA.Posts[i].Disliked = sqlitecommands.GetUserScoreOnPost(data.DB, env.MAINPAGEDATA.Posts[i].PostId, userId, "posts_dislikes")
		} else {
			env.MAINPAGEDATA.Posts[i].Liked = false
			env.MAINPAGEDATA.Posts[i].Disliked = false
		}
	}

	utility.CheckForCookies(data.DB, r, w)
	if err := templ.Execute(w, env.MAINPAGEDATA); err != nil {
		panic(err)
	}

	env.MAINPAGEDATA.GenerateAlert("", "")
}
