package model

import (
	"context"
	"time"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/ent/captcha"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func CreateCaptcha(ctx context.Context, client *ent.Client) func(text string) (*ent.Captcha, error) {
	return func(text string) (*ent.Captcha, error) {
		c, err := client.Captcha.Create().SetText(text).Save(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "failint creating captcha[%s]", text)
		}
		return c, nil
	}
}

func GetCaptchaByKey(ctx context.Context, client *ent.Client) func(key uuid.UUID) (*ent.Captcha, error) {
	return func(key uuid.UUID) (*ent.Captcha, error) {
		c, err := client.Captcha.
			Query().
			Where(captcha.Key(key)).
			Where(captcha.ExpiredTimeGTE(time.Now())).
			Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "failing getting captcha[%s]", key.String())
		}
		return c, nil
	}
}

func DeleteCapachaByKey(ctx context.Context, client *ent.Client) func(key uuid.UUID) error {
	return func(key uuid.UUID) error {
		_, err := client.Captcha.Delete().Where(captcha.Key(key)).Exec(ctx)
		if err != nil {
			return errors.Wrapf(err, "failing deleting captcha[%s]", key)
		}
		return nil
	}
}

func DeleteExpiredCaptcha(ctx context.Context, client *ent.Client) func() error {
	return func() error {
		_, err := client.Captcha.Delete().Where(captcha.ExpiredTimeLT(time.Now())).Exec(ctx)
		if err != nil {
			return errors.Wrapf(err, "failing deleting expired captchas")
		}
		return nil
	}
}
