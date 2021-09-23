package controller

import (
	"net/http"

	"github.com/YanxinTang/blog/server/e"
	"github.com/YanxinTang/blog/server/model"
	"github.com/gin-gonic/gin"
)

func GetSetting(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.Error(e.New(http.StatusBadRequest, "设置项名称不能为空"))
		return
	}
	rawValue := model.GetSetting(key)
	var value interface{}
	switch rawValue {
	case "null":
		value = nil
	case "1":
		value = true
	case "0":
		value = false
	default:
		value = rawValue
	}

	c.JSON(http.StatusOK, gin.H{
		"key":   key,
		"value": value,
	})
}

func SetSetting(c *gin.Context) {
	var setting model.Setting
	if err := c.Bind(&setting); err != nil {
		return
	}

	if err := model.SetSetting(setting.Key, setting.Value); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}
