package pages

import (
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"forum/internal/utility"
	"net/http"
)

type Main struct {
}

func (data Main) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	utility.CheckForCookies(r, w)

	if r.URL.Path != "/" {
		if utility.RedirectToPostPage(r.URL.Path) {
			LoadPostPage(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	env.MAINPAGEDATA.Posts = utility.CollectAllPostsData()
	env.MAINPAGEDATA.Categories = sqlitecommands.GetAllCategories()
	utility.CheckForCookies(r, w)
	env.MAINPAGEDATA.Username = ""

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

	if err := env.TEMPLATES["main.html"].Execute(w, env.MAINPAGEDATA); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	env.MAINPAGEDATA.GenerateAlert("", "")
}
