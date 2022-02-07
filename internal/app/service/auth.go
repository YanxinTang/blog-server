package service

import (
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"go.uber.org/zap"
)

func InitUserAndSetting(user model.User) e.ApiError {
	tx, err := model.DB.Begin(ctx)
	if err != nil {
		log.Error("InitUserAndSetting 事务创建失败", zap.Error(err))
		return e.ERROR_BEGIN_TX
	}
	defer tx.Rollback(ctx)

	if err := model.CreateUserTx(tx, user.Username, user.Email, user.RawPassword); err != nil {
		log.Error("创建初始用户失败", zap.Error(err))
		return e.ERROR_POPULATE_USER
	}
	if err := model.SetSettingTx(tx, "signupEnable", "0"); err != nil {
		log.Error("关闭注册功能失败", zap.Error(err))
		return e.ERROR_UPDAET_SETTING
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Error("InitUserAndSetting 事务提交失败", zap.Error(err))
		return e.ERROR_COMMIT_TX
	}
	return nil
}
