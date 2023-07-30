package cmd

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/config"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/file_parser"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/repo"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/type_transfer"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/xfile"
)

var configPath = "config.yaml"

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use: "run",

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

		clickhouse := config.Clickhouse{}
		cnf.UnmarshalKey("clickhouse", &clickhouse)
		config.C.Clickhouse = clickhouse

		file := config.File{}
		cnf.UnmarshalKey("file", &file)
		config.C.File = file
		fmt.Println("config load [ok]")
	},
	Run: func(cmd *cobra.Command, args []string) {
		reader, err := xfile.NewFileReader(config.C.File.Path)
		if err != nil {
			fmt.Println(err)
		}
		fileParser, exist := file_parser.DefaultParserController.GetParser("file")
		if !exist {
			fileParser = &file_parser.FileParser{}
		}
		clickhouseConf := repo.ClickhouseRepoConfig{
			Host:     config.C.Clickhouse.Host,
			Port:     config.C.Clickhouse.Port,
			DB:       config.C.Clickhouse.DB,
			User:     config.C.Clickhouse.Credentials.User,
			Password: config.C.Clickhouse.Credentials.Password,
		}
		repo, err := repo.NewClickhouseRepo(clickhouseConf)
		if err != nil {
			fmt.Println(err)
			return
		}

		columns, index, types := config.C.Clickhouse.BuildColumns()

		fileParser.SetFormat(index)
		var finish bool

		for {
			time.Sleep(10 * time.Millisecond)
			vals := [][]interface{}{}
			for i := 0; i < 43; i++ {
				b, err := reader.ReadLine()
				if err != nil {
					if err == io.EOF {
						fmt.Println("finish")
						finish = true
						break
					}
					fmt.Println("err", err)
					continue
				}
				res := fileParser.Parse(string(b))

				val := []interface{}{}
				for idx, _ := range columns {
					v := res[index[idx]]
					typ := types[idx]
					switch typ {
					case "time":
						vtime := type_transfer.String2Time(v)
						val = append(val, vtime)
					case "int32":
						vint := type_transfer.String2Int32(v)
						val = append(val, vint)
					case "int64":
						vint := type_transfer.String2Int64(v)
						val = append(val, vint)
					default:
						val = append(val, v)
					}

				}
				vals = append(vals, val)
			}
			err = repo.BatchInsert(context.Background(), "engine_log1", columns, vals, true)
			if err != nil {
				fmt.Println(err)
			}

			if finish {
				return
			}
		}

	},
}
