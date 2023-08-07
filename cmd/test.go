package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylitianrui/file-clickhouse-exporter/internal/preprocessing"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/config"
)

// var configPath = "config.yaml"

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use: "test",

	Short: "print the version number of file-clickhouse-exporter",
	Long:  `All software has versions. This is the version number of file-clickhouse-exporter`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 初始化配置
		cnf := config.Default()
		cnf.SetCnfFileName(configPath)
		err := cnf.Load()
		if err != nil {
			panic(err)
		}

		clickhouse := config.ClickHouse{}
		cnf.UnmarshalKey("clickhouse", &clickhouse)
		config.C.ClickHouse = clickhouse

		setting := config.Setting{}
		cnf.UnmarshalKey("setting", &setting)
		config.C.Setting = setting
		fmt.Println("config load [ok]")
	},
	Run: func(cmd *cobra.Command, args []string) {
		preprocessor := preprocessing.NewPreprocessor()
		columns := []string{}
		for _, v := range config.C.ClickHouse.Columns {
			columns = append(columns, v)

		}
		preprocessor.SetReadColumns(columns)
		// for k, v := range config.C.ClickHouse.Preprocessing {
		// 	preprocessor.SetProcessorLogic(k, v)
		// }
		preprocessor.SetProcessing(config.C.ClickHouse.Preprocessing)
		preprocessor.Load()

	},
}