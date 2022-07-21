package controller

import (
	"net/http"

	"github.com/YanxinTang/blog-server/internal/app/service"
	"github.com/gin-gonic/gin"
)

type SetSettingReqBody []service.SetSettingsPair

func SetSettings(c *gin.Context) {
	var setSettingReqBody SetSettingReqBody
	if err := c.BindJSON(&setSettingReqBody); err != nil {
		return
	}

	settings, err := service.SetSettings(setSettingReqBody)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, settings)
}

type GetSettingReqQuery struct {
	Keys []string `form:"keys[]" binding:"required"`
}

func PublicGetSettings(c *gin.Context) {
	var query GetSettingReqQuery
	if err := c.BindQuery(&query); err != nil {
		return
	}
	setting, err := service.GetPublicSettings(query.Keys)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, setting)
}

func GetSettings(c *gin.Context) {
	var query GetSettingReqQuery
	if err := c.BindQuery(&query); err != nil {
		return
	}
	setting, err := service.GetSettings(query.Keys)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, setting)
}
