package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dir string

func init() {
	migrateCmd.PersistentFlags().StringVarP(&dir, "dir", "d", "migrations", "migrations dir")
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	Long:  "migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}
