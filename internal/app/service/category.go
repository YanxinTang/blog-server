package service

import (
	"fmt"
	"net/http"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/YanxinTang/blog-server/utils"
)

func DeleteCategory(categoryID int) error {
	return utils.WithTx(common.Context, common.Client, func(tx *ent.Tx) error {
		count, err := model.GetCategoryArticlesCount(common.Context, tx.Client())(categoryID, model.StatusNil)
		if err != nil {
			return e.New(http.StatusInternalServerError, "获取分类所属文章数量失败，无法删除分类")
		}

		if count > 0 {
			return e.New(http.StatusBadRequest, fmt.Sprintf("由于分类下有 %d 篇文章，无法删除分类", count))
		}
		if err := model.DeleteCategory(common.Context, tx.Client())(categoryID); err != nil {
			return e.New(http.StatusInternalServerError, "分类删除失败")
		}
		return nil
	})

}
