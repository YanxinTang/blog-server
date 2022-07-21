package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/YanxinTang/blog-server/config"
	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/google/uuid"
	"github.com/steambap/captcha"
	"go.uber.org/zap"
)

func CaptchaStoreGC() {
	ticker := time.NewTicker(config.CaptchaExpiration)
	for {
		<-ticker.C
		if err := model.DeleteExpiredCaptcha(common.Context, common.Client)(); err != nil {
			log.Warn("failing collecting captcha garbage", zap.Error(err))
			return
		}
		log.Info("collect captcha garbage successfully")
	}
}

func NewCaptcha() (*captcha.Data, *ent.Captcha, error) {
	cdata, err := captcha.New(150, 50)
	if err != nil {
		log.Warn("failed to generate captcha", zap.Error(err))
		return nil, nil, e.New(http.StatusInternalServerError, "生成验证码失败")
	}
	c, err := model.CreateCaptcha(common.Context, common.Client)(cdata.Text)
	if err != nil {
		log.Warn("failed to store captcha into postgre", zap.Error(err))
		return nil, nil, e.New(http.StatusInternalServerError, "验证码存储失败")
	}
	return cdata, c, nil
}

func GetCaptcha(key string) error {
	kuuid, err := uuid.Parse(key)
	if err != nil {
		log.Warn("failing parse captcha key", zap.String("key", key), zap.Error(err))
		return e.ERROR_BAD_REQUEST
	}
	_, err = model.GetCaptchaByKey(common.Context, common.Client)(kuuid)
	if err != nil {
		return e.New(http.StatusBadRequest, "验证码未找到")
	}
	return nil
}

func VerifyCaptcha(key, text string) error {
	kuuid, err := uuid.Parse(key)
	if err != nil {
		log.Warn("failing parsing captcha key", zap.String("key", key), zap.Error(err))
		return e.ERROR_BAD_REQUEST
	}

	captcha, err := model.GetCaptchaByKey(common.Context, common.Client)(kuuid)
	if err != nil {
		log.Warn("failing getting captcha", zap.String("key", key))
		return e.New(http.StatusBadRequest, "验证码输入有误")
	}

	if !strings.EqualFold(captcha.Text, text) {
		log.Warn("failing verifing captcha", zap.String("requestCaptchaText", text), zap.String("storeCaptchaText", captcha.Text))
		return e.New(http.StatusBadRequest, "验证码输入有误")
	}
	return nil
}
