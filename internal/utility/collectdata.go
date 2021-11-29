package utility

import (
	"database/sql"
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"log"
	"net/http"
	"time"
)

func CollectPostData(db *sql.DB) env.Post {
	var post env.Post
	var date time.Time
	var body string
	var authorId int

	_ = db.QueryRow("SELECT * FROM posts WHERE id = ?", env.POSTID).Scan(&post.PostId, &authorId, &post.Title, &body, &date, &post.Likes, &post.DisLikes, &post.Comments)

	post.Created = date.Format("January 02, 2006 at 15:04")
	post.Author = sqlitecommands.GetUserNameFromTable(db, authorId)
	post.Categories = sqlitecommands.GetPostCategoriesFromTable(db, post.PostId)
	post.Body = DivideBodyIntoParagraphs(body)

	return post
}

func CollectAllPostsData(db *sql.DB) []env.Post {
	var result []env.Post

	rows, err := db.Query("SELECT * FROM posts")
	CheckErr(err)

	for rows.Next() {
		var post env.Post
		var date time.Time
		var body string
		var authorId int

		err := rows.Scan(&post.PostId, &authorId, &post.Title, &body, &date, &post.Likes, &post.DisLikes, &post.Comments)
		if err != nil {
			log.Fatal(err)
		}

		post.Created = date.Format("January 02, 2006 at 15:04")
		post.Author = sqlitecommands.GetUserNameFromTable(db, authorId)
		post.Categories = sqlitecommands.GetPostCategoriesFromTable(db, post.PostId)
		post.Image.Name, post.Image.Container, post.Image.Type = sqlitecommands.GetImageDataFromTable(db, post.PostId)
		post.Body = DivideBodyIntoParagraphs(body)

		result = append(result, post)
	}

	return result
}

func CollectAllPostComments(db *sql.DB, postId int, w http.ResponseWriter, r *http.Request) []env.Comment {
	var result []env.Comment
	var userId int = -1

	rows, err := db.Query("SELECT * FROM posts_comments WHERE post_id = ?", postId)
	CheckErr(err)

	if CheckForCookies(db, r, w) {
		userId = sqlitecommands.GetUserIdByCookies(db, r, w)
	}

	for rows.Next() {
		var comment env.Comment
		var body string
		var date time.Time
		var id, postId, authorId int

		err := rows.Scan(&id, &postId, &authorId, &comment.Likes, &comment.Dislikes, &date, &body)
		if err != nil {
			log.Fatal(err)
		}

		comment.Id = id
		comment.Author = sqlitecommands.GetUserNameFromTable(db, authorId)
		comment.Created = date.Format("January 02, 2006 at 15:04")
		comment.Body = DivideBodyIntoParagraphs(body)

		if userId != -1 {
			if sqlitecommands.GetUserScoreOnComment(db, comment.Id, userId, "comment_likes") {
				comment.Liked = "liked"
			} else if sqlitecommands.GetUserScoreOnComment(db, comment.Id, userId, "comment_dislikes") {
				comment.Disliked = "disliked"
			}
		}

		result = append(result, comment)
	}

	return result
}
