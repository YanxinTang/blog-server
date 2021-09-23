package model

type Setting struct {
	BaseModel
	Key   string      `db:"key" json:"key"`
	Value interface{} `db:"value" json:"value"`
}

func GetSetting(key string) (value string) {
	DB.Get(&value, "SELECT `value` FROM `setting` WHERE `key` = ?", key)
	return
}

func SetSetting(key string, value interface{}) error {
	if _, err := DB.Exec(
		"INSERT INTO `setting` (`key`, `value`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `value` = ?",
		key,
		value,
		value,
	); err != nil {
		return err
	}
	return nil
}
