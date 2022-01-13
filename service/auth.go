package service

import (
	"log"

	"github.com/YanxinTang/blog-server/e"
	"github.com/YanxinTang/blog-server/model"
)

func InitUserAndSetting(user model.User) e.ApiError {
	tx, err := model.DB.Begin(ctx)
	if err != nil {
		log.Println(err)
		return e.ERROR_BEGIN_TX
	}
	defer tx.Rollback(ctx)

	if err := model.CreateUserTx(tx, user.Username, user.Email, user.RawPassword); err != nil {
		log.Println(err)
		return e.ERROR_POPULATE_USER
	}
	if err := model.SetSettingTx(tx, "signupEnable", "0"); err != nil {
		log.Println(err)
		return e.ERROR_UPDAET_SETTING
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Println(err)
		return e.ERROR_COMMIT_TX
	}
	return nil
}
