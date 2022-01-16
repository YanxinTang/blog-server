package model

import "github.com/georgysavva/scany/pgxscan"

type Storage struct {
	BaseModel
	Name      string `db:"name" json:"name"`
	SecretID  string `db:"secret_id" json:"secretID"`
	SecretKey string `db:"secret_key" json:"secretKey"`
	Token     string `db:"token" json:"token"`
	Region    string `db:"region" json:"region"`
	Endpoint  string `db:"endpoint" json:"endpoint"`
	Bucket    string `db:"bucket" json:"bucket"`
	Usage     int64  `db:"usage" json:"usage"`
	Capacity  int64  `db:"capacity" json:"capacity"`
}

func GetStorages() ([]Storage, error) {
	storages := []Storage{}
	err := pgxscan.Select(ctx, db, &storages, "SELECT * FROM storage")
	return storages, err
}

func GetStorage(ID uint64) (Storage, error) {
	var storage Storage
	err := pgxscan.Get(ctx, db, &storage, "SELECT * FROM storage WHERE id = $1", ID)
	return storage, err
}

func CreateStorage(storage Storage) (Storage, error) {
	sql := `
	INSERT INTO storage
	(name, secret_id, secret_key, token, region, endpoint, bucket, capacity)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING *
	`
	err := pgxscan.Get(
		ctx, db, &storage,
		sql,
		storage.Name,
		storage.SecretID,
		storage.SecretKey,
		storage.Token,
		storage.Region,
		storage.Endpoint,
		storage.Bucket,
		storage.Capacity,
	)
	return storage, err
}

func UpdateStorage(storage Storage) error {
	sql := `
	UPDATE storage 
	SET (name, secret_id, secret_key, token, region, endpoint, bucket, capacity)
	= ($1, $2, $3, $4, $5, $6, $7, $8)
	WHERE id = $9
	`
	_, err := db.Exec(ctx, sql,
		storage.Name,
		storage.SecretID,
		storage.SecretKey,
		storage.Token,
		storage.Region,
		storage.Endpoint,
		storage.Bucket,
		storage.Capacity,
		storage.ID,
	)
	return err
}

func DeleteStorage(storageID uint64) error {
	_, err := db.Exec(ctx, "DELETE FROM storage WHERE id = $1", storageID)
	return err
}
