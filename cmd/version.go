package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version holds the application version
var Version string

func init() {
	Version = "v0.1.1-dev0"
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
