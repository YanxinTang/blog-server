package controller

import (
	"net/http"

	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/gin-gonic/gin"
)

type apiGetSettingQuery struct {
	Key string `form:"key" binding:"required"`
}

func GetSetting(c *gin.Context) {
	var query apiGetSettingQuery
	if err := c.BindQuery(&query); err != nil {
		return
	}
	setting, err := model.GetSetting(query.Key)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":   query.Key,
		"value": setting.Value,
	})
}

type apiSetSettingBody struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func SetSetting(c *gin.Context) {
	var apiSetSetting apiSetSettingBody
	if err := c.Bind(&apiSetSetting); err != nil {
		return
	}

	if err := model.SetSetting(apiSetSetting.Key, apiSetSetting.Value); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}
