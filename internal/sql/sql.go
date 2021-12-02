package sqlitecommands

import (
	"database/sql"
	"forum/internal/env"
	"log"
	"net/http"
)

func UpdateUsersTable(sessionToken string, userName string, email string, password string, date string, role int, ip string) {
	stmt, err := env.DB.Prepare("INSERT INTO users(username, email, password, date, role, ip) values(?,?,?,?,?,?)")
	CheckErr(err)

	result, err := stmt.Exec(userName, email, password, date, role, ip)
	CheckErr(err)

	uid, _ := result.LastInsertId()
	UpdateSessionToken(sessionToken, int(uid))
}

func UpdatePostsTable(authorId int, postTitle string, postContent string, date string, imagePath string, postCategories []string) {
	stmt, err := env.DB.Prepare("INSERT INTO posts(author_id, title, body, created, likes, dislikes, comments) values(?,?,?,?,?,?,?)")
	CheckErr(err)

	result, err := stmt.Exec(authorId, postTitle, postContent, date, 0, 0, 0)
	CheckErr(err)

	postId, _ := result.LastInsertId()
	UpdatePostsPicturesTable(int(postId), imagePath)
	UpdatePostsCategoriesTable(int(postId), postCategories)
}

func UpdatePostsPicturesTable(postId int, imagePath string) {
	stmt, err := env.DB.Prepare("INSERT INTO posts_images(post_id, image_path) values(?,?)")
	CheckErr(err)

	_, err = stmt.Exec(postId, imagePath)
	CheckErr(err)
}

func UpdatePostsCategoriesTable(postId int, categories []string) {
	for i := 0; i < len(categories); i++ {
		stmt, err := env.DB.Prepare("INSERT INTO posts_categories(post_id, category) values(?,?)")
		CheckErr(err)

		_, err = stmt.Exec(postId, categories[i])
		CheckErr(err)
	}

	UpdateCategoriesTable(categories)
}

func UpdateCategoriesTable(categories []string) {
	for i := 0; i < len(categories); i++ {
		if !FindSameCategory(categories[i]) {
			stmt, err := env.DB.Prepare("INSERT INTO categories(category) values(?)")
			CheckErr(err)

			_, err = stmt.Exec(categories[i])
			CheckErr(err)
		}
	}
}

func UpdatePostsCommentsTable(postId int, authorId int, created string, body string) {
	var value int

	stmt, err := env.DB.Prepare("INSERT INTO posts_comments(post_id, author_id, likes, dislikes, created, body) values(?,?,?,?,?,?)")
	CheckErr(err)

	_, err = stmt.Exec(postId, authorId, 0, 0, created, body)
	CheckErr(err)

	counterStmt := "SELECT COUNT (*) FROM posts_comments WHERE post_id = ?"
	_ = env.DB.QueryRow(counterStmt, postId).Scan(&value)
	CheckErr(err)

	updateStmt, err := env.DB.Prepare("UPDATE posts SET comments = ? WHERE id = ?")
	CheckErr(err)

	_, err = updateStmt.Exec(value, postId)
	CheckErr(err)
}

func UpdateCommentRatingTable(commentId int, authorId int, updateType string, table string) {
	var stmt *sql.Stmt
	var err error

	if updateType == "add" {
		stmt, err = env.DB.Prepare("INSERT INTO " + table + " (comment_id, author_id) values(?,?)")
	} else if updateType == "remove" {
		stmt, err = env.DB.Prepare("DELETE FROM " + table + " WHERE comment_id = ? AND author_id = ?")
	}

	CheckErr(err)

	stmt.Exec(commentId, authorId)
}

func UpdatePostRatingTable(postId int, authorId int, updateType string, table string) {
	var stmt *sql.Stmt
	var err error

	if updateType == "add" {
		stmt, err = env.DB.Prepare("INSERT INTO " + table + "(post_id, author_id) values(?,?)")
	} else if updateType == "remove" {
		stmt, err = env.DB.Prepare("DELETE FROM " + table + " WHERE post_id = ? AND author_id = ?")
	}

	CheckErr(err)

	stmt.Exec(postId, authorId)
}

func UpdateSessionToken(sessionToken string, uid int) {
	stmt, err := env.DB.Prepare("UPDATE users SET session_token = ? WHERE uid = ?")
	CheckErr(err)
	_, err = stmt.Exec(sessionToken, uid)
	CheckErr(err)
}

func GetPostsIdGap() (int, int) {
	var first, last int

	if CheckExistence("posts") {
		return -1, -1
	}

	selectStmt := "SELECT id FROM posts ORDER BY id DESC LIMIT 1"
	err := env.DB.QueryRow(selectStmt).Scan(&last)
	CheckErr(err)

	selectStmt = "SELECT id FROM posts LIMIT 1"
	err = env.DB.QueryRow(selectStmt).Scan(&first)
	CheckErr(err)

	return first, last
}

func GetUserName(id int) string {
	var result string

	selectStmt := "SELECT username FROM users WHERE uid = ?"
	err := env.DB.QueryRow(selectStmt, id).Scan(&result)
	CheckErr(err)

	return result
}

func GetPostCategories(id int) []string {
	var result []string

	rows, err := env.DB.Query("SELECT category FROM posts_categories WHERE post_id = ?", id)
	CheckErr(err)

	for rows.Next() {
		var value string

		err := rows.Scan(&value)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, value)
	}

	return result
}

func GetPostRatingCounter(postId int) (int, int) {
	var likes, dislikes int

	counterStmt := "SELECT COUNT (*) FROM posts_likes WHERE post_id = ?"
	env.DB.QueryRow(counterStmt, postId).Scan(&likes)

	counterStmt = "SELECT COUNT (*) FROM posts_dislikes WHERE post_id = ?"
	env.DB.QueryRow(counterStmt, postId).Scan(&dislikes)

	return likes, dislikes
}

func GetCommentRatingCounter(commentId int) (int, int) {
	var likes, dislikes int

	counterStmt := "SELECT COUNT (*) FROM comment_likes WHERE comment_id = ?"
	env.DB.QueryRow(counterStmt, commentId).Scan(&likes)

	counterStmt = "SELECT COUNT (*) FROM comment_dislikes WHERE comment_id = ?"
	env.DB.QueryRow(counterStmt, commentId).Scan(&dislikes)

	return likes, dislikes
}

func GetPostCommentCounter(postId int) int {
	var counter int

	counterStmt := "SELECT COUNT (*) FROM posts_comments WHERE post_id = ?"
	env.DB.QueryRow(counterStmt, postId).Scan(&counter)

	return counter
}

func GetImagePath(id int) string {
	var imagePath string

	selectStmt := "SELECT image_path FROM posts_images WHERE post_id = ?"
	err := env.DB.QueryRow(selectStmt, id).Scan(&imagePath)
	CheckErr(err)

	return imagePath
}

func GetIntData(tableName string, columnName string, id int, idColumnName string) int {
	var result int

	return result
}

func GetUserScoreOnPost(postId int, authorId int, tableName string) bool {
	var value int

	sqlStmt := "SELECT COUNT (*) FROM " + tableName + " WHERE post_id = ? and author_id = ?"
	_ = env.DB.QueryRow(sqlStmt, postId, authorId).Scan(&value)

	return value == 1
}

func GetUserScoreOnComment(commentId int, authorId int, tableName string) bool {
	var value int

	sqlStmt := "SELECT COUNT (*) FROM " + tableName + " WHERE comment_id = ? and author_id = ?"
	_ = env.DB.QueryRow(sqlStmt, commentId, authorId).Scan(&value)

	return value == 1
}

func GetUserIdByCookies(r *http.Request, w http.ResponseWriter) int {
	c, err := r.Cookie("session_token")

	if err == nil {
		uid, checkResult := CheckDataExistence(c.Value, "session_token")

		if checkResult {
			return uid
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	return -1
}

func GetAllCategories() []string {
	var result []string

	sqlStmt := "SELECT category FROM categories"
	rows, err := env.DB.Query(sqlStmt)
	CheckErr(err)

	for rows.Next() {
		var category string

		err := rows.Scan(&category)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, category)
	}

	return result
}

func FindSameCategory(category string) bool {
	sqlStmt := "SELECT category FROM categories WHERE category = ?"
	err := env.DB.QueryRow(sqlStmt, category).Scan(&category)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return false
	}
	return true
}

func CheckDataExistence(REGDATA string, dataType string) (int, bool) {
	var uid int
	sqlStmt := "SELECT " + dataType + ", uid FROM users WHERE " + dataType + " = ?"
	err := env.DB.QueryRow(sqlStmt, REGDATA).Scan(&REGDATA, &uid)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return -1, false
	}
	return uid, true
}

func CheckExistence(tableName string) bool {
	sqlStmt := "SELECT * FROM " + tableName
	err := env.DB.QueryRow(sqlStmt).Scan()

	return err == sql.ErrNoRows
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
