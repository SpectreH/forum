package sqlitecommands

import (
	"database/sql"
	"fmt"
	"log"
)

func UpdateUsersTable(db *sql.DB, sessionToken string, userName string, email string, password string, date string, role int, ip string) {
	stmt, err := db.Prepare("INSERT INTO users(username, email, password, date, role, ip) values(?,?,?,?,?,?)")
	CheckErr(err)

	result, err := stmt.Exec(userName, email, password, date, role, ip)
	CheckErr(err)

	uid, _ := result.LastInsertId()
	UpdateSessionToken(db, sessionToken, int(uid))
}

func UpdatePostsTable(db *sql.DB, authorId int, postTitle string, postContent string, date string, imageName string, imageContainer []byte, imageType string, postCategories []string) {
	stmt, err := db.Prepare("INSERT INTO posts(author_id, title, body, created, likes, dislikes, comments) values(?,?,?,?,?,?,?)")
	CheckErr(err)

	result, err := stmt.Exec(authorId, postTitle, postContent, date, 0, 0, 0)
	CheckErr(err)

	postId, _ := result.LastInsertId()
	UpdatePostsPicturesTable(db, int(postId), imageName, imageContainer, imageType)
	UpdatePostsCategoriesTable(db, int(postId), postCategories)
}

func UpdatePostsPicturesTable(db *sql.DB, postId int, imageName string, imageContainer []byte, imageType string) {
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
	stmt, err := db.Prepare("INSERT INTO posts_comments(post_id, author_id, created, body) values(?,?,?,?)")
	CheckErr(err)

	_, err = stmt.Exec(postId, authorId, created, body)
	CheckErr(err)
}

func UpdateRatingsTable(db *sql.DB, postId int, authorId int, updateType string, rating string) {
	var stmt *sql.Stmt
	var err error

	if rating != "posts_likes" && rating != "posts_dislikes" {
		return
	}

	if updateType == "add" {
		stmt, err = db.Prepare("INSERT INTO " + rating + "(post_id, author_id) values(?,?)")
	} else {
		stmt, err = db.Prepare("DELETE FROM " + rating + " WHERE post_id = ? AND author_id = ?")
	}
	CheckErr(err)

	_, err = stmt.Exec(postId, authorId)
	CheckErr(err)
}

func UpdateSessionToken(db *sql.DB, sessionToken string, uid int) {
	stmt, err := db.Prepare("UPDATE users SET session_token = ? WHERE uid = ?")
	CheckErr(err)
	_, err = stmt.Exec(sessionToken, uid)
	CheckErr(err)
}

func UpdatePostsData(db *sql.DB, postId int, column string, updateType string) {
	updateStmt, err := db.Prepare("UPDATE posts SET " + column + " = ? WHERE id = ?")
	CheckErr(err)

	var value int
	selectStmt := "SELECT " + column + " FROM posts WHERE id = ?"
	err = db.QueryRow(selectStmt, postId).Scan(&value)
	if err != nil {
		fmt.Println("error")
		return
	}

	if updateType == "sum" {
		value += 1
	} else {
		value -= 1
	}

	_, err = updateStmt.Exec(value, postId)
	CheckErr(err)
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

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
