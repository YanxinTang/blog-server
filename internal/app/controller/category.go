package controller

import (
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/app/service"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/YanxinTang/blog-server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetCategories(c *gin.Context) {
	categories, err := model.GetCategories(common.Context, common.Client)()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, categories)
}

type CreateCategoryReqBody struct {
	Name string `json:"name"`
}

func CareteCategory(c *gin.Context) {
	var createCategoryReqBody CreateCategoryReqBody
	if err := c.BindJSON(&createCategoryReqBody); err != nil {
		return
	}

	category, err := model.CreateCategory(common.Context, common.Client)(createCategoryReqBody.Name)
	if err != nil {
		log.Warn(
			"failint creating category",
			zap.String("name", createCategoryReqBody.Name),
			zap.Error(err),
		)
		c.Error(e.New(http.StatusBadRequest, "分类名称有误，请重新输入"))
		return
	}
	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	categoryID, err := utils.GetID(c, "categoryID")
	if err != nil {
		return
	}

	if err := service.DeleteCategory(categoryID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

type UpdateCategoryReqbody struct {
	Name string `json:"name"`
}

func UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("categoryID"))
	if err != nil {
		log.Warn("failing converting categoryID", zap.String("categoryID", c.Param("categoryID")))
		c.Error(e.ERROR_BAD_REQUEST)
		return
	}
	var updateCategoryReqBody UpdateCategoryReqbody
	if err := c.BindJSON(&updateCategoryReqBody); err != nil {
		return
	}

	var uci model.UpdateCategoryInput
	uci.ID = categoryID
	uci.Name = updateCategoryReqBody.Name

	category, err := model.UpdateCategory(common.Context, common.Client)(uci)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, category)
}
