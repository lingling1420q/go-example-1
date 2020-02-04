package main

import (
	"fmt"
	"os"

	"github.com/gunsluo/go-example/migrate/build"
	"github.com/gunsluo/go-example/migrate/migrate"
	"github.com/spf13/cobra"
)

var verbose bool
var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "cmd",
		Short: "Example Command",
		Long:  "Top level command for Example Command.",
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(
		migrate.Cmd,
		build.Cmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}