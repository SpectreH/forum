package pages

import (
	"forum/internal/env"
	"forum/internal/utility"
	"net/http"
)

type Logout struct {
}

func (data Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if utility.CheckForCookies(r, w) {
		c := http.Cookie{
			Name:   "session_token",
			MaxAge: -1}
		http.SetCookie(w, &c)

		env.MAINPAGEDATA.GenerateAlert("You have successfully logged out!", "Logout")
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
