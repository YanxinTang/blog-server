package cmd

import (
	"fmt"
	"os"

	"github.com/YanxinTang/blog-server/internal/app/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "The backends of simple blog",
	Long:  "The backends of simple blog",
	Run: func(cmd *cobra.Command, args []string) {
		svr := server.New()
		svr.Run(":8000")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
