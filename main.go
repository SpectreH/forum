package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"forum/internal/env"
	"forum/internal/pages"
)

func main() {
	env.InitEnv()

	env.DB, _ = sql.Open("sqlite3", "./db/forum.db")
	defer env.DB.Close()

	os.Mkdir("images", 0700)

	http.Handle("/", pages.Main{})
	http.Handle("/login", pages.Login{})
	http.Handle("/registration", pages.Registration{})
	http.Handle("/new", pages.New{})
	http.Handle("/logout", pages.Logout{})
	http.Handle("/account", pages.Account{})

	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	js := http.FileServer(http.Dir("js"))
	http.Handle("/js/", http.StripPrefix("/js/", js))

	images := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", images))

	fmt.Println("Server is listening on port 8000...")
	if http.ListenAndServe(":8000", nil) != nil {
		log.Fatalf("%v - Internal Server Error", http.StatusInternalServerError)
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
