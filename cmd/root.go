package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/YanxinTang/blog/server/config"
	"github.com/YanxinTang/blog/server/model"
	"github.com/YanxinTang/blog/server/router"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "The backends of simple blog",
	Long:  "The backends of simple blog",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.ParseConfig()
		if err != nil {
			log.Fatal(err)
		}
		pool, err := config.GetDBConnectionPool(conf.Postgres)
		if err != nil {
			log.Fatal(err)
		}
		store, err := config.GetCookieStore(*conf)
		if err != nil {
			log.Fatal(err)
		}
		model.Setup(pool)
		engine := router.SetupRouter(store)
		engine.Run(":8000")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
