package controller

import (
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog/server/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	categories, err := model.GetCategories()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, categories)
}

func CareteCategory(c *gin.Context) {
	var category model.Category
	if err := c.ShouldBind(&category); err != nil {
		c.Error(err)
		return
	}

	session := sessions.Default(c)
	userID := session.Get("userID").(uint64)

	category, err := model.CreateCategory(userID, category)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, category)
}

func UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("categoryID"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	var category model.Category
	if err := c.ShouldBind(&category); err != nil {
		c.Error(err)
		return
	}

	category.ID = categoryID

	err = model.UpdateCategory(category)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
