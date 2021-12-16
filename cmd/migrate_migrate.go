package cmd

import (
	"fmt"

	"github.com/YanxinTang/blog/server/config"
	"github.com/spf13/cobra"
)

var version uint

func init() {
	migrateMigrateCmd.Flags().UintVarP(&version, "version", "v", 0, "step to migrate to")
	migrateMigrateCmd.MarkFlagRequired("version")
	migrateCmd.AddCommand(migrateMigrateCmd)
}

var migrateMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate to specify version",
	Long:  "Migrate looks at the currently active migration version, then migrates either up or down to the specified version.",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := config.GetDefaultMigrate()
		if err != nil {
			exitWithError(err)
		}
		defer migrater.Close()
		if err := migrater.Migrate(version); err != nil {
			exitWithError(err)
		}
		fmt.Println("Successfully!")
	},
}
