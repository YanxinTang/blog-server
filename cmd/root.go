package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/YanxinTang/blog-server/config"
	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/app/server"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "The backends of simple blog",
	Long:  "The backends of simple blog",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.ParseConfig()
		if err != nil {
			log.Fatal("failed to parse config", zap.Error(err))
		}

		client, err := config.GetEntClient(conf)
		if err != nil {
			log.Fatal("failed opening connection to sqlite: %v", zap.Error(err))
		}
		defer client.Close()
		common.SetupEntClient(client)
		if err := model.AutoMigrate(context.Background(), client); err != nil {
			log.Fatal("failed creating schema resources: %v", zap.Error(err))
		}
		if err := model.Seed(context.Background(), client); err != nil {
			log.Warn("failing seeding", zap.Error(err))
		}

		store, err := config.GetCookieStore(conf)
		if err != nil {
			log.Fatal("failed to get cookie store", zap.Error(err))
		}

		svr := server.New(store)
		svr.Run(":8000")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
