package cmd

import (
	"fmt"

	"github.com/YanxinTang/blog-server/config"
	"github.com/spf13/cobra"
)

var step int

func init() {
	migrateStepsCmd.Flags().IntVarP(&step, "step", "s", 0, "step to migrate to")
	migrateStepsCmd.MarkFlagRequired("step")
	migrateCmd.AddCommand(migrateStepsCmd)
}

var migrateStepsCmd = &cobra.Command{
	Use:   "steps",
	Short: "Steps looks at the currently active migration version",
	Long:  "Steps looks at the currently active migration version. It will migrate up if n > 0, and down if n < 0.",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, err := config.GetDefaultMigrate()
		if err != nil {
			exitWithError(err)
		}
		defer migrater.Close()
		if err := migrater.Steps(step); err != nil {
			exitWithError(err)
		}
		fmt.Println("Successfully!")
	},
}
