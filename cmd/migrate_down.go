package cmd

import (
	"fmt"

	"github.com/YanxinTang/blog/server/config"
	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.AddCommand(migrateDownCmd)
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "migrate down",
	Long:  "migrate down",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := config.GetDefaultMigrate()
		if err != nil {
			exitWithError(err)
		}
		defer migrater.Close()
		if err := migrater.Down(); err != nil {
			exitWithError(err)
		}
		fmt.Println("applying all down migrations")
	},
}
