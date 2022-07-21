package common

import (
	"context"

	"github.com/YanxinTang/blog-server/ent"
)

var Client *ent.Client
var Context = context.Background()

func SetupEntClient(entClient *ent.Client) {
	Client = entClient
}
