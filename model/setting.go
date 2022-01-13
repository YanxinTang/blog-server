package model

import (
	"github.com/georgysavva/scany/pgxscan"
)

type Setting struct {
	BaseModel
	Key   string `db:"key"`
	Value string `db:"value"`
}

func GetSetting(key string) (Setting, error) {
	var setting Setting
	err := pgxscan.Get(ctx, db, &setting, "SELECT value FROM setting WHERE key = $1", key)
	return setting, err
}

func SetSetting(key string, value string) error {
	return SetSettingTx(db, key, value)
}

func SetSettingTx(tx Executor, key string, value string) error {
	if _, err := tx.Exec(
		ctx,
		"INSERT INTO setting (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value",
		key,
		value,
	); err != nil {
		return err
	}
	return nil
}
