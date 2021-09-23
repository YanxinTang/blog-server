package model

import (
	"database/sql"
	"log"
)

type Comment struct {
	BaseModel
	ArticleID       uint64 `db:"article_id" json:"articleID"`
	ParentCommentID uint64 `db:"parent_comment_id" json:"parentCommentID"`
	Username        string `db:"username" json:"username"`
	Content         string `db:"content" json:"content" binding:"required"`
}

func getComment(commentID uint64) (Comment, error) {
	row := DB.QueryRowx("SELECT * FROM comment WHERE id = ?", commentID)
	var comment Comment
	err := row.StructScan(&comment)
	if err != nil {
		return comment, err
	}
	return comment, nil
}

func CreateComment(comment Comment) (Comment, error) {
	res, err := DB.Exec(
		"INSERT comment SET article_id = ?, username = ?, content = ?",
		comment.ArticleID,
		comment.Username,
		comment.Content,
	)
	if err != nil {
		return comment, err
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return comment, err
	}
	lastInsertComment, err := getComment(uint64(lastInsertID))
	if err != nil {
		return comment, err
	}
	return lastInsertComment, nil
}

func DeleteComment(commentID uint64) (sql.Result, error) {
	stmt, err := DB.Prepare("DELETE FROM comment WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(commentID)
}

func GetArticleComments(articleID uint64) ([]Comment, error) {
	rows, err := DB.Queryx("SELECT * FROM comment WHERE article_id = ?", articleID)
	if err != nil {
		return nil, err
	}
	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.StructScan(&comment)
		if err != nil {
			log.Println(err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// CommentsCount returns count of comment
func CommentsCount() uint64 {
	var count uint64
	DB.QueryRow("SELECT COUNT(*) FROM comment").Scan(&count)
	return count
}
