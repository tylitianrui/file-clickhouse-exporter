package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/tylitianrui/file-clickhouse-exporter/internal/preprocessing"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/aggregation"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/config"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/file_parser"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/repo"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/xfile"
)

var configPath = "config.yaml"
var mu sync.Mutex

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use: "run",

	Short: "read file and insert db",
	Long:  `read file and insert clickhouse`,
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
		clickHouseConfig := config.C.ClickHouse
		dbConfig := repo.ClickhouseRepoConfig{
			Host:     clickHouseConfig.Host,
			Port:     clickHouseConfig.Port,
			DB:       clickHouseConfig.DB,
			User:     clickHouseConfig.Credentials.User,
			Password: clickHouseConfig.Credentials.Password,
		}
		db, err := repo.NewClickhouseRepo(dbConfig)
		if err != nil {
			fmt.Println(err)
			return
		}
		preprocessor := preprocessing.NewPreprocessor()
		preprocessor.SetColumns(config.C.ClickHouse.Columns)
		preprocessor.SetPreprocessingConfig(config.C.ClickHouse.Preprocessing)
		err = preprocessor.LoadConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
		var reader xfile.XReader
		switch config.C.Setting.Mode {
		case "follower":
			reader, err = xfile.NewAppendingFileAppendReader(config.C.Setting.FilePath)
			if err != nil {
				fmt.Println(err)
				return
			}

		default:
			reader, err = xfile.NewStaticFileReader(config.C.Setting.FilePath)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		// todo cancel and signal
		ctx, _ := context.WithCancel(context.Background())
		lines := reader.ReadLines(ctx)
		idx := preprocessor.GetIndexOfFile()
		parser, _ := file_parser.DefaultParserController.GetParser("file")
		parser.SetFormat(idx)
		columns := preprocessor.Columns()
		vals := [][]interface{}{}
		var num int
		ticker := time.NewTicker(time.Duration(config.C.Setting.Interval) * time.Millisecond)
		go func() {
			for {
				select {
				case <-ticker.C:
					if len(vals) > 0 {
						mu.Lock()
						err := db.BatchInsert(context.Background(), clickHouseConfig.Table, columns, vals, false)
						fmt.Println(err)
						num = 0
						vals = vals[:0]
						mu.Unlock()
					}
				default:
					if num >= config.C.Setting.MaxlineEveryRead && len(vals) > 0 {
						mu.Lock()
						err := db.BatchInsert(context.Background(), clickHouseConfig.Table, columns, vals, false)
						fmt.Println(err)
						num = 0
						vals = vals[:0]
						mu.Unlock()
					}
					time.Sleep(100 * time.Millisecond)
					// time.Sleep(time.Duration(config.C.Setting.Interval) * time.Millisecond)
				}
			}
		}()
		for line := range lines {
			result := preprocessing.NewResult()
			data := parser.Parse(string(line.Line()))
			result.SetRaw(data)
			aggregations := preprocessor.Do(data)
			result.SetAggregation(aggregations)
			dynamic, exist := config.C.ClickHouse.Preprocessing[preprocessing.PreprocessorDynamic]
			if exist {
				dynamicResult := map[string]string{}
				for k, d := range dynamic {
					switch d {
					case "gen_uuid()":
						dynamicResult[k] = aggregation.GenUUID()
					default:

					}
				}
				result.SetDynamic(dynamicResult)

			}
			static, exist := config.C.ClickHouse.Preprocessing[preprocessing.PreprocessorStatic]
			if exist {
				result.SetStatic(static)
			}
			resIdx := preprocessor.ResultIdx()
			res := result.Result(resIdx)
			val := preprocessor.ColumnsResult(res)
			mu.Lock()
			vals = append(vals, val)
			num++
			mu.Unlock()

		}

	},
}
