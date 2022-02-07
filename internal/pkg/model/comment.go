package model

import (
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
)

type Comment struct {
	BaseModel
	ArticleID uint64 `db:"article_id" json:"articleID"`
	ParentID  uint64 `db:"parent_id" json:"parentID"`
	Username  string `db:"username" json:"username"`
	Content   string `db:"content" json:"content" binding:"required"`
}

func CreateComment(comment Comment) (Comment, error) {
	err := pgxscan.Get(
		ctx, db, &comment,
		"INSERT INTO comment (article_id, username, content) VALUES ($1, $2, $3) RETURNING *",
		comment.ArticleID,
		comment.Username,
		comment.Content,
	)
	return comment, err
}

func DeleteComment(commentID uint64) (pgconn.CommandTag, error) {
	return db.Exec(ctx, "DELETE FROM comment WHERE id = $1", commentID)
}

func GetArticleComments(articleID uint64) ([]Comment, error) {
	comments := []Comment{}
	err := pgxscan.Select(
		ctx,
		db,
		&comments,
		"SELECT * FROM comment WHERE article_id = $1",
		articleID,
	)
	return comments, err
}

// CommentsCount returns count of comment
func CommentsCount() (uint64, error) {
	var count uint64
	err := pgxscan.Get(ctx, db, &count, "SELECT COUNT(*) FROM comment")
	return count, err
}
