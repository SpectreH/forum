package pages

import (
	"forum/internal/env"
	"forum/internal/utility"
	"net/http"
)

type Account struct {
}

func (data Account) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !utility.CheckForCookies(r, w) {
		utility.RedirectToMainPage(r, w, "You are not logged in!", "Fail_NotLoggedIn")
		return
	}

	if err := env.TEMPLATES["error.html"].Execute(w, env.MAINPAGEDATA); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
