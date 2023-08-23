package cmd

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "file-clickhouse-exporter",
	Short: "file-clickhouse-exporter",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		go func() {
			http.ListenAndServe("0.0.0.0:8080", nil)
		}()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		// handle  error
		return err
	}
	return nil
}
