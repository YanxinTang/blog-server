package main

import (
	"github.com/YanxinTang/blog-server/cmd"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
)

func main() {
	defer log.Sync()
	cmd.Execute()
}
