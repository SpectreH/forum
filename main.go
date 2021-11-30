package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"forum/internal/pages"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/forum.db")
	CheckErr(err)

	defer db.Close()

	http.Handle("/", pages.Main{DB: db})
	http.Handle("/login", pages.Login{DB: db})
	http.Handle("/registration", pages.Registration{DB: db})
	http.Handle("/new", pages.New{DB: db})
	http.Handle("/logout", pages.Logout{DB: db})

	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	js := http.FileServer(http.Dir("js"))
	http.Handle("/js/", http.StripPrefix("/js/", js))

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
