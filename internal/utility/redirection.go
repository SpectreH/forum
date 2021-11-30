package utility

import (
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"net/http"
	"strconv"
	"strings"
)

func RedirectToMainPage(r *http.Request, w http.ResponseWriter, message string, alertType string) {
	env.MAINPAGEDATA.GenerateAlert(message, alertType)

	http.Redirect(w, r, "/", http.StatusFound)
}

func RedirectToPostPage(URL string) bool {
	first, last := sqlitecommands.GetPostsIdGap()
	if first == -1 && last == -1 {
		return false
	}

	postId := strings.Trim(URL, "/")
	number, err := strconv.Atoi(postId)
	if err == nil {
		if number <= last && first <= number {
			env.POSTID = number
			return true
		}
	}
	return false
}
