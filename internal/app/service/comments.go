package service

import (
	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/YanxinTang/blog-server/internal/pkg/page"
)

func GetPublishedArticleComments(articleID int, p *page.Pagination) ([]*ent.Comment, *page.Pagination, error) {
	comments, p, err := model.GetArticleComments(common.Context, common.Client)(articleID, model.StatusPublished, p)
	if err != nil {
		return nil, nil, e.ERROR_RESOURCE_NOT_FOUND
	}
	return comments, p, nil
}
