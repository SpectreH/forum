package pages

import (
	"fmt"
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"html/template"
	"net/http"
)

type Main struct {
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

	env.MAINPAGEDATA.Posts = utility.CollectAllPostsData()
	env.MAINPAGEDATA.Categories = sqlitecommands.GetAllCategories()

	utility.CheckForCookies(r, w)
	env.MAINPAGEDATA.Username = "-1"

	var userId int
	if env.MAINPAGEDATA.LoggedIn {
		userId = sqlitecommands.GetUserIdByCookies(r, w)
		env.MAINPAGEDATA.Username = sqlitecommands.GetUserName(userId)
	}

	for i := 0; i < len(env.MAINPAGEDATA.Posts); i++ {
		if env.MAINPAGEDATA.LoggedIn {
			env.MAINPAGEDATA.Posts[i].Liked = sqlitecommands.GetUserScoreOnPost(env.MAINPAGEDATA.Posts[i].PostId, userId, "posts_likes")
			env.MAINPAGEDATA.Posts[i].Disliked = sqlitecommands.GetUserScoreOnPost(env.MAINPAGEDATA.Posts[i].PostId, userId, "posts_dislikes")
		} else {
			env.MAINPAGEDATA.Posts[i].Liked = false
			env.MAINPAGEDATA.Posts[i].Disliked = false
		}
	}

	if err := templ.Execute(w, env.MAINPAGEDATA); err != nil {
		panic(err)
	}

	env.MAINPAGEDATA.GenerateAlert("", "")
}
