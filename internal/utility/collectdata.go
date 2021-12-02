package utility

import (
	"forum/internal/env"
	sqlitecommands "forum/internal/sql"
	"log"
	"net/http"
	"time"
)

func CollectPostData() env.Post {
	var post env.Post
	var date time.Time
	var body string
	var authorId int

	_ = env.DB.QueryRow("SELECT * FROM posts WHERE id = ?", env.POSTID).Scan(&post.PostId, &authorId, &post.Title, &body, &date, &post.Likes, &post.DisLikes, &post.Comments)

	post.Created = date.Format("January 02, 2006 at 15:04")
	post.Author = sqlitecommands.GetUserName(authorId)
	post.Categories = sqlitecommands.GetPostCategories(post.PostId)
	post.Body = DivideBodyIntoParagraphs(body)
	post.Likes, post.DisLikes = sqlitecommands.GetPostRatingCounter(env.POSTID)
	post.Comments = sqlitecommands.GetPostCommentCounter(env.POSTID)

	return post
}

func CollectAllPostsData() []env.Post {
	var result []env.Post

	rows, err := env.DB.Query("SELECT * FROM posts")
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
		post.Author = sqlitecommands.GetUserName(authorId)
		post.Categories = sqlitecommands.GetPostCategories(post.PostId)
		post.ImagePath = sqlitecommands.GetImagePath(post.PostId)
		post.Body = DivideBodyIntoParagraphs(body)
		post.Likes, post.DisLikes = sqlitecommands.GetPostRatingCounter(post.PostId)
		post.Comments = sqlitecommands.GetPostCommentCounter(env.POSTID)

		result = append(result, post)
	}

	return result
}

func CollectAllPostComments(postId int, w http.ResponseWriter, r *http.Request) []env.Comment {
	var result []env.Comment
	var userId int = -1

	rows, err := env.DB.Query("SELECT * FROM posts_comments WHERE post_id = ?", postId)
	CheckErr(err)

	if CheckForCookies(r, w) {
		userId = sqlitecommands.GetUserIdByCookies(r, w)
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
		comment.Author = sqlitecommands.GetUserName(authorId)
		comment.Created = date.Format("January 02, 2006 at 15:04")
		comment.Body = DivideBodyIntoParagraphs(body)
		comment.Likes, comment.Dislikes = sqlitecommands.GetCommentRatingCounter(comment.Id)

		if userId != -1 {
			if sqlitecommands.GetUserScoreOnComment(comment.Id, userId, "comment_likes") {
				comment.Liked = true
			} else if sqlitecommands.GetUserScoreOnComment(comment.Id, userId, "comment_dislikes") {
				comment.Disliked = true
			}
		}

		result = append(result, comment)
	}

	return result
}
