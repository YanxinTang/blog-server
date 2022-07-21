package service

import (
	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/YanxinTang/blog-server/utils"
)

func InitUserAndSetting(user ent.User) error {
	return utils.WithTx(common.Context, common.Client, func(tx *ent.Tx) error {
		_, err := model.CreateUser(common.Context, tx.Client())(user.Username, user.Email, user.Password)
		return err
	})
}
