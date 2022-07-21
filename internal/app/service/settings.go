package service

import (
	"net/http"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/YanxinTang/blog-server/utils"
	"go.uber.org/zap"
)

type SetSettingsPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func SetSettings(pairs []SetSettingsPair) ([]*ent.Setting, error) {
	settings := make([]*ent.Setting, 0, len(pairs))
	err := utils.WithTx(common.Context, common.Client, func(tx *ent.Tx) error {
		for _, pair := range pairs {
			if pair.Key == "" {
				return e.New(http.StatusBadRequest, "设置项键不能为空")
			}
			if pair.Value == "" {
				return e.New(http.StatusBadRequest, "设置项值不能为空")
			}
			setting, err := model.SetSetting(common.Context, tx.Client())(pair.Key, pair.Value)
			if err != nil {
				return e.New(http.StatusBadRequest, "设置保存失败")
			}
			settings = append(settings, setting)
		}
		return nil
	})
	return settings, err
}

func GetSettings(keys []string) ([]*ent.Setting, error) {
	settings := make([]*ent.Setting, len(keys))
	for i, key := range keys {
		setting, err := model.GetSetting(common.Context, common.Client)(key)
		if err != nil {
			if !ent.IsNotFound(err) {
				return nil, e.ERROR_RESOURCE_NOT_FOUND
			}
			settings[i] = &ent.Setting{Key: key}
			log.Warn("failing getting setting, return emtpy alternative", zap.String("key", key))
		} else {
			settings[i] = setting
		}
	}
	return settings, nil
}

var publicSettingKeys []string = []string{
	"signupEnable",
}

func GetPublicSettings(keys []string) ([]*ent.Setting, error) {
	for _, key := range keys {
		if !utils.ArrayStringIncludes(publicSettingKeys, key) {
			return nil, e.ERROR_RESOURCE_NOT_FOUND
		}
	}
	return GetSettings(keys)
}
