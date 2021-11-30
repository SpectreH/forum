package utility

import (
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func CreateSessionToken(w http.ResponseWriter) string {
	sessionToken := uuid.NewV4().String()
	env.MAINPAGEDATA.LoggedIn = true

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(1200 * time.Second),
	})
	return sessionToken
}

func CheckForCookies(r *http.Request, w http.ResponseWriter) bool {
	c, err := r.Cookie("session_token")
	if err == nil {
		_, checkResult := sqlitecommands.CheckDataExistence(c.Value, "session_token")

		if checkResult {
			env.MAINPAGEDATA.LoggedIn = true
			return true
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	env.MAINPAGEDATA.LoggedIn = false
	return false
}
