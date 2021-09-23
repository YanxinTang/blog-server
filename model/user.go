package model

import (
	"encoding/base64"
	"log"

	"github.com/YanxinTang/blog/server/e"
	"github.com/YanxinTang/blog/server/utils"
)

type User struct {
	BaseModel
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	Salt     string `josn:"-" db:"salt"`
}

func getUserByUsername(username string) (User, error) {
	res := DB.QueryRowx("SELECT * FROM `user` WHERE `username` = ?", username)
	var user User
	err := res.StructScan(&user)
	return user, err
}

func Authentication(username, password string) (User, error) {
	user, err := getUserByUsername(username)
	if err != nil {
		log.Printf("getUserByUsername(%s): %s", username, err)
		return user, e.ERROR_INVALID_AUTH
	}

	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		log.Printf("base64.StdEncoding.DecodeString(%s): %s", user.Salt, err)
		return user, e.ERROR_INVALID_AUTH
	}

	if !utils.DoPasswordsMatch(user.Password, password, salt) {
		log.Printf("base64.StdEncoding.DecodeString(%s, %s): %s", user.Password, password, err)
		return user, e.ERROR_INVALID_AUTH
	}

	return user, nil
}

func CreateUser(username, email, password string) error {
	salt := utils.GenerateRandomSalt()
	hashPassword := utils.HashPassword(password, salt)
	_, err := DB.Exec(
		"INSERT INTO `user` (`username`, `email`, `password`, `salt` ) VALUES (?, ?, ?, ?)",
		username,
		email,
		hashPassword,
		base64.StdEncoding.EncodeToString(salt),
	)
	return err
}
