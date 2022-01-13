package cmd

import (
	"fmt"

	"github.com/YanxinTang/blog-server/config"
	"github.com/spf13/cobra"
)

var forceVersion int

func init() {
	migrateForceCmd.Flags().IntVarP(&forceVersion, "version", "v", 0, "version to set forcely")
	migrateForceCmd.MarkFlagRequired("version")
	migrateCmd.AddCommand(migrateForceCmd)
}

var migrateForceCmd = &cobra.Command{
	Use:   "force",
	Short: "Force sets a migration version",
	Long:  "Force sets a migration version. It does not check any currently active version in database. It resets the dirty state to false.",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := config.GetDefaultMigrate()
		if err != nil {
			exitWithError(err)
		}
		defer migrater.Close()
		if err := migrater.Force(forceVersion); err != nil {
			exitWithError(err)
		}
		fmt.Println("Successfully!")
	},
}
