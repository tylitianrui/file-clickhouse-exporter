package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylitianrui/file-clickhouse-exporter/internal"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "print the version number of file-clickhouse-exporter",
	Long:    `All software has versions. This is the version number of file-clickhouse-exporter`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(internal.VERSION)
	},
}
