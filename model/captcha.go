package model

import (
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type Captcha struct {
	BaseModel
	Key  string `db:"key"`
	Text string `db:"text"`
}

func CreateCaptcha(captcha Captcha) (Captcha, error) {
	err := pgxscan.Get(
		ctx, db, &captcha,
		"INSERT INTO captcha (key, text) VALUES ($1, $2) RETURNING id, created_at, updated_at, key, text",
		captcha.Key,
		captcha.Text,
	)
	return captcha, err
}

func GetCaptchaByKey(key string, expiration time.Duration) (Captcha, error) {
	expiredat := time.Now().Add(-expiration)
	var captcha Captcha
	err := pgxscan.Get(
		ctx,
		db,
		&captcha,
		"SELECT id, key, text, created_at, updated_at FROM captcha WHERE key = $1 AND created_at > $2",
		key,
		expiredat,
	)
	return captcha, err
}

func DeleteCapachaByKey(key string) error {
	_, err := db.Exec(ctx, "DELETE FROM captcha WHERE key = $1", key)
	return err
}

func DeleteExpiredCaptcha(expiration time.Duration) error {
	expiredat := time.Now().Add(-expiration)
	_, err := db.Exec(ctx, "DELETE FROM captcha WHERE created_at < $1", expiredat)
	return err
}
