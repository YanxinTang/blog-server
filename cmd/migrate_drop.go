package cmd

import (
	"fmt"

	"github.com/YanxinTang/blog-server/config"
	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.AddCommand(migrateDropCmd)
}

var migrateDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop deletes everything in the database",
	Long:  "Drop deletes everything in the database",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := config.GetDefaultMigrate()
		if err != nil {
			exitWithError(err)
		}
		defer migrater.Close()
		if err := migrater.Drop(); err != nil {
			exitWithError(err)
		}
		fmt.Println("applying all up migrations")
	},
}
