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

		clickhouse := config.ClickHouse{}
		cnf.UnmarshalKey("clickhouse", &clickhouse)
		config.C.ClickHouse = clickhouse

		setting := config.Setting{}
		cnf.UnmarshalKey("setting", &setting)
		config.C.Setting = setting
		fmt.Println("config load [ok]")
	},
	Run: func(cmd *cobra.Command, args []string) {
		reader, err := xfile.NewFileReader(config.C.Setting.FilePath)
		if err != nil {
			fmt.Println(err)
		}
		fileParser, exist := file_parser.DefaultParserController.GetParser("file")
		if !exist {
			fileParser = &file_parser.FileParser{}
		}
		clickhouseConf := repo.ClickhouseRepoConfig{
			Host:     config.C.ClickHouse.Host,
			Port:     config.C.ClickHouse.Port,
			DB:       config.C.ClickHouse.DB,
			User:     config.C.ClickHouse.Credentials.User,
			Password: config.C.ClickHouse.Credentials.Password,
		}
		repo, err := repo.NewClickhouseRepo(clickhouseConf)
		if err != nil {
			fmt.Println(err)
			return
		}

		preprocessing, err := config.C.ClickHouse.BuildColumns()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fileParser.SetFormat(preprocessing.Index)
		var finish bool

		for {
			time.Sleep(time.Duration(config.C.Setting.Interval) * time.Millisecond)
			vals := [][]interface{}{}
			for i := 0; i < config.C.Setting.MaxlineEveryRead; i++ {
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
				for idx, _ := range preprocessing.Columns {
					v := res[preprocessing.Index[idx]]
					split := preprocessing.Split[idx]
					f := split[0]
					t := split[1]
					if f != 0 && t != -1 {
						v = v[f:t]
					} else if f != 0 {
						v = v[f:]
					} else if t != -1 {
						v = v[:t]
					}
					typ := preprocessing.Types[idx]
					switch typ {
					case "time":
						vtime := type_transfer.String2Time(v)
						val = append(val, vtime)
					case "time_utc":
						vtime := type_transfer.String2TimeUTC(v)
						val = append(val, vtime)
					case "int32":
						vint := type_transfer.String2Int32(v)
						val = append(val, vint)
					case "int64":
						vint := type_transfer.String2Int64(v)
						val = append(val, vint)
					case "string":
						fallthrough
					default:
						val = append(val, v)
					}

				}
				vals = append(vals, val)
			}
			err = repo.BatchInsert(context.Background(), config.C.ClickHouse.Table, preprocessing.Columns, vals, true)
			if err != nil {
				fmt.Println("[err] err:", err)
				fmt.Println(fmt.Sprintf("[failure] saved %d rows ", len(vals)))
			} else {
				fmt.Println(fmt.Sprintf("[success] saved %d rows ", len(vals)))
			}

			if finish {
				return
			}
		}

	},
}
