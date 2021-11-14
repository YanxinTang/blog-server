package model

import "github.com/georgysavva/scany/pgxscan"

type Setting struct {
	BaseModel
	Key   string      `db:"key" json:"key"`
	Value interface{} `db:"value" json:"value"`
}

func GetSetting(key string) (value string) {
	pgxscan.Get(ctx, db, &value, "SELECT value FROM setting WHERE key = $1", key)
	return
}

func SetSetting(key string, value interface{}) error {
	if _, err := db.Exec(
		ctx,
		"INSERT INTO setting (key, value, type) VALUES ('$1', $2, 0) ON CONFLICT (key) DO UPDATE SET value = excluded.value",
		key,
		value,
	); err != nil {
		return err
	}
	return nil
}
