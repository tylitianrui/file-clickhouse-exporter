package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "file-clickhouse-exporter",
	Short: "file-clickhouse-exporter",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		// handle  error
		return err
	}
	return nil
}
