package controller

import (
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog/server/model"
	"github.com/gin-gonic/gin"
)

func DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("categoryID"), 10, 64)
	if err != nil {
		c.Error(err).SetType(http.StatusNotFound)
		return
	}
	if err := model.DeleteCategory(categoryID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
