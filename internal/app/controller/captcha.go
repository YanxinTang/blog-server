package controller

import (
	"bytes"
	"encoding/base64"
	"net/http"

	"github.com/YanxinTang/blog-server/internal/app/service"
	"github.com/gin-gonic/gin"
)

func GetCapacha(c *gin.Context) {
	cdata, cmodel, apierr := service.NewCaptcha()
	if apierr != nil {
		c.Error(apierr)
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, 128))
	if err := cdata.WriteImage(buf); err != nil {
		c.Error(err)
		return
	}
	cbase64Data := base64.StdEncoding.EncodeToString(buf.Bytes())

	c.JSON(http.StatusOK, gin.H{
		"key":  cmodel.Key,
		"data": cbase64Data,
	})
}

type VerifyCaptchaReq struct {
	Key  string `json:"captchaKey" binding:"required"`
	Text string `json:"captchaText" binding:"required"`
}
