package cmd

import (
	"fmt"
	"log"

	"github.com/YanxinTang/blog/server/config"
	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.AddCommand(migrateVersionCmd)
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version returns the currently active migration version",
	Long:  "Version returns the currently active migration version",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := config.GetDefaultMigrate()
		if err != nil {
			log.Fatal(err)
		}
		defer migrater.Close()
		version, dirty, err := migrater.Version()
		if err != nil {
			exitWithError(err)
		}
		fmt.Printf("Version\tDirty\n%d\t%t\n", version, dirty)
	},
}
