package model

import (
	"context"

	"github.com/YanxinTang/blog-server/ent"
)

func AutoMigrate(ctx context.Context, client *ent.Client) error {
	// Run the auto migration tool.
	return client.Schema.Create(ctx)
}
