package env

import (
	"database/sql"
	"text/template"
)

type Image struct {
	Name      string
	Type      string
	Container string
}

type Post struct {
	PostId     int
	Author     string
	Title      string
	Body       []string
	Created    string
	Likes      int
	DisLikes   int
	Comments   int
	Categories []string
	Image      Image
	Liked      bool
	Disliked   bool
}

type Comment struct {
	Id       int
	Author   string
	Created  string
	Body     []string
	Likes    int
	Dislikes int
	Liked    bool
	Disliked bool
}

type MainPage struct {
	Message    string
	AlertType  string
	Posts      []Post
	Categories []string
	LoggedIn   bool
	Username   string
}

type PostPage struct {
	Message   string
	AlertType string
	Post      Post
	Comments  []Comment
	LoggedIn  bool
}

var MAINPAGEDATA MainPage
var POSTID int
var DB *sql.DB
var TEMPLATES map[string]*template.Template

func InitEnv() {
	MAINPAGEDATA = MainPage{
		Message:   "",
		AlertType: "",
		Posts:     nil,
		LoggedIn:  false,
	}

	TEMPLATES = ParseTemplates()
}

func (data MainPage) GenerateAlert(message string, alertType string) {
	MAINPAGEDATA.Message = message
	MAINPAGEDATA.AlertType = alertType
}
