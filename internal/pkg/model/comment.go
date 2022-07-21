package model

import (
	"context"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/ent/article"
	"github.com/YanxinTang/blog-server/ent/comment"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/internal/pkg/page"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CreateCommentInput struct {
	ArticleID int
	ParentID  int
	Username  string
	Content   string
}

func CreateComment(ctx context.Context, client *ent.Client) func(CreateCommentInput) (*ent.Comment, error) {
	return func(input CreateCommentInput) (*ent.Comment, error) {
		exec := client.Comment.
			Create().
			SetArticleID(input.ArticleID).
			SetUsername(input.Username).
			SetContent(input.Content)
		if input.ParentID > 0 {
			exec = exec.SetParentID(input.ParentID)
		}
		c, err := exec.Save(ctx)
		if err != nil {
			log.Warn("failing creating comment", zap.Error(err))
			return nil, err
		}
		return c, nil
	}
}

func DeleteComment(ctx context.Context, client *ent.Client) func(id int) error {
	return func(id int) error {
		err := client.Comment.DeleteOneID(id).Exec(ctx)
		if err != nil {
			return errors.Wrapf(err, "failing deleting comment[%d]", id)
		}
		return nil
	}
}

func DeleteArticleComments(ctx context.Context, client *ent.Client) func(articleID int) error {
	return func(articleID int) error {
		n, err := client.Comment.Delete().Where(comment.HasArticleWith(article.ID(articleID))).Exec(ctx)
		if err != nil {
			log.Warn("failing deleting comments of article", zap.Int("articleID", articleID), zap.Error(err))
			return err
		}
		log.Info("delete comments of articles", zap.Int("count", n))
		return nil
	}
}

func GetComment(ctx context.Context, client *ent.Client) func(id int) (*ent.Comment, error) {
	return func(id int) (*ent.Comment, error) {
		c, err := client.Comment.Get(ctx, id)
		if err != nil {
			return nil, errors.Wrapf(err, "failing getting comment[%d]", id)
		}
		return c, nil
	}
}

func GetArticleComments(ctx context.Context, client *ent.Client) func(articleID int, status ArticleStatus, p *page.Pagination) ([]*ent.Comment, *page.Pagination, error) {
	return func(articleID int, status ArticleStatus, p *page.Pagination) ([]*ent.Comment, *page.Pagination, error) {
		query := client.Comment.Query().Where(comment.HasArticleWith(article.ID(articleID)))
		comments, err := query.Clone().
			Order(ent.Asc(comment.FieldID)).
			Offset(p.Offset()).Limit(p.Limit()).
			All(ctx)
		if err != nil {
			log.Warn("failing getting comments of article", zap.Int("articleID", articleID), zap.Int8("status", int8(status)), zap.Error(err))
			return nil, nil, err
		}
		total, err := query.Clone().Count(ctx)
		if err != nil {
			log.Warn("failing getting comments count of article", zap.Int("articleID", articleID), zap.Int8("status", int8(status)), zap.Error(err))
			return nil, nil, err
		}
		p.Total = total
		return comments, p, nil
	}
}

func CommentsCount(ctx context.Context, client *ent.Client) func() (int, error) {
	return func() (int, error) {
		count, err := client.Comment.Query().Count(ctx)
		if err != nil {
			return 0, errors.Wrap(err, "failing getting count of comment")
		}
		return count, nil
	}
}
