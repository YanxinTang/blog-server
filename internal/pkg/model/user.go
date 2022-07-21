package model

import (
	"context"
	"encoding/base64"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/ent/user"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func getUserByUsername(ctx context.Context, client *ent.Client) func(username string) (*ent.User, error) {
	return func(username string) (*ent.User, error) {
		user, err := client.User.
			Query().
			Where(user.Username(username)).
			Only(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed querying user")
		}
		log.Info("user returned", zap.Int("userID", user.ID), zap.String("username", user.Username))
		return user, nil
	}
}

func Authentication(ctx context.Context, client *ent.Client) func(username, password string) (*ent.User, error) {
	return func(username, password string) (*ent.User, error) {
		user, err := getUserByUsername(ctx, client)(username)
		if err != nil {
			log.Error(
				"failing getting user",
				zap.String("username", username),
				zap.Error(err),
			)
			return nil, errors.Wrapf(err, "failing getting user[%s]", username)
		}

		if !utils.DoPasswordsMatch(user.Password, password, user.Salt) {
			log.Warn(
				"password mismatch",
				zap.Int("userID", user.ID),
				zap.String("user password", user.Password),
				zap.String("request password", password),
				zap.String("salt", base64.StdEncoding.EncodeToString(user.Salt)),
				zap.Error(err),
			)
			return nil, errors.New("password mismatch")
		}

		return user, nil
	}
}

func CreateUser(ctx context.Context, client *ent.Client) func(username, email, password string) (*ent.User, error) {
	return func(username, email, password string) (*ent.User, error) {
		salt := utils.GenerateRandomSalt()
		hashPassword := utils.HashPassword(password, salt)

		user, err := client.User.
			Create().
			SetUsername(username).
			SetEmail(email).
			SetPassword(hashPassword).
			SetSalt(salt).
			Save(ctx)
		return user, err
	}
}
