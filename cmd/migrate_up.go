package cmd

import (
	"fmt"

	"github.com/YanxinTang/blog-server/config"
	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "migrate up",
	Long:  "migrate up",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := config.GetDefaultMigrate()
		if err != nil {
			exitWithError(err)
		}
		defer migrater.Close()
		if err := migrater.Up(); err != nil {
			exitWithError(err)
		}
		fmt.Println("applying all up migrations")
	},
}
