package model

import (
	"context"

	"github.com/YanxinTang/blog-server/ent"
)

func Seed(ctx context.Context, client *ent.Client) error {
	if err := client.Setting.Create().SetKey("signupEnable").SetValue("1").OnConflictColumns("key").DoNothing().Exec(ctx); err != nil {
		return err
	}
	return nil
}
