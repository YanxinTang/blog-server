package model

import (
	"context"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/ent/setting"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"go.uber.org/zap"
)

func GetSetting(ctx context.Context, client *ent.Client) func(key string) (*ent.Setting, error) {
	return func(key string) (*ent.Setting, error) {
		s, err := client.Setting.
			Query().
			Where(setting.Key(key)).
			Only(ctx)
		if err != nil {
			return nil, err
		}
		log.Info("setting returned", zap.String("key", s.Key), zap.String("value", s.Value))
		return s, nil
	}
}

func SetSetting(ctx context.Context, client *ent.Client) func(key, value string) (*ent.Setting, error) {
	return func(key, value string) (*ent.Setting, error) {
		if id, err := client.Setting.Query().Where(setting.Key(key)).OnlyID(ctx); err != nil {
			// only log error
			log.Warn("failing quering the setting", zap.String("key", key), zap.Error(err))
		} else if id > 0 {
			// update setting with new value
			s, err := client.Setting.UpdateOneID(id).SetValue(value).Save(ctx)
			if err != nil {
				log.Warn("failing updatig the setting", zap.Int("id", id), zap.Error(err))
				return nil, err
			}
			return s, nil
		}

		// the setting with `key` is not found, just create it
		s, err := client.Setting.
			Create().
			SetKey(key).
			SetValue(value).
			Save(ctx)
		if err != nil {
			log.Warn("failing creating setting", zap.String("key", key), zap.String("value", value), zap.Error(err))
			return nil, err
		}
		return s, nil
	}
}
