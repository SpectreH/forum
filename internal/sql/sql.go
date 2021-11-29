package sqlitecommands

import (
	"database/sql"
	"log"
	"net/http"
)

func UpdateUsersTable(db *sql.DB, sessionToken string, userName string, email string, password string, date string, role int, ip string) {
	stmt, err := db.Prepare("INSERT INTO users(username, email, password, date, role, ip) values(?,?,?,?,?,?)")
	CheckErr(err)

	result, err := stmt.Exec(userName, email, password, date, role, ip)
	CheckErr(err)

	uid, _ := result.LastInsertId()
	UpdateSessionToken(db, sessionToken, int(uid))
}

func UpdatePostsTable(db *sql.DB, authorId int, postTitle string, postContent string, date string, imageName string, imageContainer string, imageType string, postCategories []string) {
	stmt, err := db.Prepare("INSERT INTO posts(author_id, title, body, created, likes, dislikes, comments) values(?,?,?,?,?,?,?)")
	CheckErr(err)

	result, err := stmt.Exec(authorId, postTitle, postContent, date, 0, 0, 0)
	CheckErr(err)

	postId, _ := result.LastInsertId()
	UpdatePostsPicturesTable(db, int(postId), imageName, imageContainer, imageType)
	UpdatePostsCategoriesTable(db, int(postId), postCategories)
}

func UpdatePostsPicturesTable(db *sql.DB, postId int, imageName string, imageContainer string, imageType string) {
	stmt, err := db.Prepare("INSERT INTO posts_images(post_id, image_name, image_container, image_type) values(?,?,?,?)")
	CheckErr(err)

	_, err = stmt.Exec(postId, imageName, imageContainer, imageType)
	CheckErr(err)
}

func UpdatePostsCategoriesTable(db *sql.DB, postId int, categories []string) {
	for i := 0; i < len(categories); i++ {
		stmt, err := db.Prepare("INSERT INTO posts_categories(post_id, category) values(?,?)")
		CheckErr(err)

		_, err = stmt.Exec(postId, categories[i])
		CheckErr(err)
	}

	UpdateCategoriesTable(db, categories)
}

func UpdateCategoriesTable(db *sql.DB, categories []string) {
	for i := 0; i < len(categories); i++ {
		if !FindSameCategory(db, categories[i]) {
			stmt, err := db.Prepare("INSERT INTO categories(category) values(?)")
			CheckErr(err)

			_, err = stmt.Exec(categories[i])
			CheckErr(err)
		}
	}
}

func UpdatePostsCommentsTable(db *sql.DB, postId int, authorId int, created string, body string) {
	var value int

	stmt, err := db.Prepare("INSERT INTO posts_comments(post_id, author_id, likes, dislikes, created, body) values(?,?,?,?,?,?)")
	CheckErr(err)

	_, err = stmt.Exec(postId, authorId, 0, 0, created, body)
	CheckErr(err)

	counterStmt := "SELECT COUNT (*) FROM posts_comments WHERE post_id = ?"
	_ = db.QueryRow(counterStmt, postId).Scan(&value)
	CheckErr(err)

	updateStmt, err := db.Prepare("UPDATE posts SET comments = ? WHERE id = ?")
	CheckErr(err)

	_, err = updateStmt.Exec(value, postId)
	CheckErr(err)
}

func UpdateCommentRatingsTable(db *sql.DB, commentId int, authorId int, updateType string, table string) {
	var stmt *sql.Stmt
	var err error
	var value int
	var row, mirrorTable, mirrorRow string

	if table != "comment_likes" && table != "comment_dislikes" {
		return
	}

	if table == "comment_likes" {
		row = "likes"
		mirrorRow = "dislikes"
		mirrorTable = "comment_dislikes"
	} else {
		row = "dislikes"
		mirrorRow = "likes"
		mirrorTable = "comment_likes"
	}

	if updateType == "add" {
		stmt, err = db.Prepare("INSERT INTO " + table + "(comment_id, author_id) values(?,?)")
		CheckErr(err)

		checkStmt, err := db.Prepare("DELETE FROM " + mirrorTable + " WHERE comment_id = ? AND author_id = ?")
		CheckErr(err)

		_, err = checkStmt.Exec(commentId, authorId)
		CheckErr(err)

		updateStmtMirrow, err := db.Prepare("UPDATE posts_comments SET " + mirrorRow + " = ? WHERE id = ?")
		CheckErr(err)

		mirrorCounterStmt := "SELECT COUNT (*) FROM " + mirrorTable + " WHERE comment_id = ?"
		_ = db.QueryRow(mirrorCounterStmt, commentId).Scan(&value)
		_, err = updateStmtMirrow.Exec(value, commentId)
		CheckErr(err)
	} else {
		stmt, err = db.Prepare("DELETE FROM " + table + " WHERE comment_id = ? AND author_id = ?")
		CheckErr(err)
	}

	_, err = stmt.Exec(commentId, authorId)
	CheckErr(err)

	updateStmt, err := db.Prepare("UPDATE posts_comments SET " + row + " = ? WHERE id = ?")
	CheckErr(err)
	counterStmt := "SELECT COUNT (*) FROM " + table + " WHERE comment_id = ?"
	_ = db.QueryRow(counterStmt, commentId).Scan(&value)
	_, err = updateStmt.Exec(value, commentId)
	CheckErr(err)
}

func UpdatePostRatingsTable(db *sql.DB, postId int, authorId int, updateType string, table string) {
	var stmt *sql.Stmt
	var err error
	var value int
	var row, mirrorTable, mirrorRow string

	if table != "posts_likes" && table != "posts_dislikes" {
		return
	}

	if table == "posts_likes" {
		row = "likes"
		mirrorRow = "dislikes"
		mirrorTable = "posts_dislikes"
	} else {
		row = "dislikes"
		mirrorRow = "likes"
		mirrorTable = "posts_likes"
	}

	if updateType == "add" {
		stmt, err = db.Prepare("INSERT INTO " + table + "(post_id, author_id) values(?,?)")
		CheckErr(err)

		checkStmt, err := db.Prepare("DELETE FROM " + mirrorTable + " WHERE post_id = ? AND author_id = ?")
		CheckErr(err)

		_, err = checkStmt.Exec(postId, authorId)
		CheckErr(err)

		updateStmtMirrow, err := db.Prepare("UPDATE posts SET " + mirrorRow + " = ? WHERE id = ?")
		CheckErr(err)

		mirrorCounterStmt := "SELECT COUNT (*) FROM " + mirrorTable + " WHERE post_id = ?"
		_ = db.QueryRow(mirrorCounterStmt, postId).Scan(&value)
		_, err = updateStmtMirrow.Exec(value, postId)
		CheckErr(err)
	} else {
		stmt, err = db.Prepare("DELETE FROM " + table + " WHERE post_id = ? AND author_id = ?")
		CheckErr(err)
	}

	_, err = stmt.Exec(postId, authorId)
	CheckErr(err)

	updateStmt, err := db.Prepare("UPDATE posts SET " + row + " = ? WHERE id = ?")
	CheckErr(err)
	counterStmt := "SELECT COUNT (*) FROM " + table + " WHERE post_id = ?"
	_ = db.QueryRow(counterStmt, postId).Scan(&value)
	_, err = updateStmt.Exec(value, postId)
	CheckErr(err)
}

func UpdateSessionToken(db *sql.DB, sessionToken string, uid int) {
	stmt, err := db.Prepare("UPDATE users SET session_token = ? WHERE uid = ?")
	CheckErr(err)
	_, err = stmt.Exec(sessionToken, uid)
	CheckErr(err)
}

func GetPostsIdGap(db *sql.DB) (int, int) {
	var first, last int

	if CheckExistence(db, "posts") {
		return -1, -1
	}

	selectStmt := "SELECT id FROM posts ORDER BY id DESC LIMIT 1"
	err := db.QueryRow(selectStmt).Scan(&last)
	CheckErr(err)

	selectStmt = "SELECT id FROM posts LIMIT 1"
	err = db.QueryRow(selectStmt).Scan(&first)
	CheckErr(err)

	return first, last
}

func GetUserNameFromTable(db *sql.DB, id int) string {
	var result string

	selectStmt := "SELECT username FROM users WHERE uid = ?"
	err := db.QueryRow(selectStmt, id).Scan(&result)
	CheckErr(err)

	return result
}

func GetPostCategoriesFromTable(db *sql.DB, id int) []string {
	var result []string

	rows, err := db.Query("SELECT category FROM posts_categories WHERE post_id = ?", id)
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

func GetImageDataFromTable(db *sql.DB, id int) (string, string, string) {
	var imageName, imageType, imageCountainer string

	selectStmt := "SELECT image_name, image_container, image_type FROM posts_images WHERE post_id = ?"
	err := db.QueryRow(selectStmt, id).Scan(&imageName, &imageCountainer, &imageType)
	CheckErr(err)

	return imageName, imageCountainer, imageType
}

func GetIntDataFromTable(db *sql.DB, tableName string, columnName string, id int, idColumnName string) int {
	var result int

	return result
}

func GetUserScoreOnPost(db *sql.DB, postId int, authorId int, tableName string) bool {
	var value int

	sqlStmt := "SELECT COUNT (*) FROM " + tableName + " WHERE post_id = ? and author_id = ?"
	_ = db.QueryRow(sqlStmt, postId, authorId).Scan(&value)

	if value == 1 {
		return true
	}

	return false
}

func GetUserScoreOnComment(db *sql.DB, commentId int, authorId int, tableName string) bool {
	var value int

	sqlStmt := "SELECT COUNT (*) FROM " + tableName + " WHERE comment_id = ? and author_id = ?"
	_ = db.QueryRow(sqlStmt, commentId, authorId).Scan(&value)

	if value == 1 {
		return true
	}

	return false
}

func GetUserIdByCookies(db *sql.DB, r *http.Request, w http.ResponseWriter) int {
	c, err := r.Cookie("session_token")

	if err == nil {
		uid, checkResult := CheckDataExistence(db, c.Value, "session_token")

		if checkResult {
			return uid
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	return -1
}

func GetAllCategoriesFromTable(db *sql.DB) []string {
	var result []string

	sqlStmt := "SELECT category FROM categories"
	rows, err := db.Query(sqlStmt)
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

func FindSameCategory(db *sql.DB, category string) bool {
	sqlStmt := "SELECT category FROM categories WHERE category = ?"
	err := db.QueryRow(sqlStmt, category).Scan(&category)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return false
	}
	return true
}

func CheckDataExistence(db *sql.DB, REGDATA string, dataType string) (int, bool) {
	var uid int
	sqlStmt := "SELECT " + dataType + ", uid FROM users WHERE " + dataType + " = ?"
	err := db.QueryRow(sqlStmt, REGDATA).Scan(&REGDATA, &uid)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return -1, false
	}
	return uid, true
}

func CheckExistence(db *sql.DB, tableName string) bool {
	sqlStmt := "SELECT * FROM " + tableName
	err := db.QueryRow(sqlStmt).Scan()

	if err == sql.ErrNoRows {
		return true
	}

	return false
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
