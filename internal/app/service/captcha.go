package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/YanxinTang/blog-server/config"
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
		log.Info("captcha store garbage collection")
		model.DeleteExpiredCaptcha(config.CaptchaExpiration)
	}
}

func NewCaptcha() (*captcha.Data, *model.Captcha, e.ApiError) {
	cdata, err := captcha.New(150, 50)
	if err != nil {
		log.Warn("failed to generate captcha", zap.Error(err))
		return nil, nil, e.New(http.StatusInternalServerError, "生成验证码失败")
	}
	uid := uuid.New()
	cmodel := model.Captcha{Key: uid.String(), Text: cdata.Text}
	cmodel, err = model.CreateCaptcha(cmodel)
	if err != nil {
		log.Warn("failed to save captcha into postgre store", zap.Error(err))
		return nil, nil, e.New(http.StatusInternalServerError, "验证码存储失败")
	}
	return cdata, &cmodel, nil
}

func GetCaptcha(key string) e.ApiError {
	_, err := model.GetCaptchaByKey(key, config.CaptchaExpiration)
	if err != nil {
		return e.New(http.StatusBadRequest, "验证码未找到")
	}
	return nil
}

func VerifyCaptcha(key, text string) e.ApiError {
	captcha, err := model.GetCaptchaByKey(key, config.CaptchaExpiration)
	if err != nil || !strings.EqualFold(captcha.Text, text) {
		return e.New(http.StatusBadRequest, "验证码输入有误")
	}
	return nil
}
