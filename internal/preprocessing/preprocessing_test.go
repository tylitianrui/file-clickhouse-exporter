package preprocessing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/file_parser"
)

func TestPreprocessor(t *testing.T) {
	a := assert.New(t)
	demo1 := "$1"
	res1 := regexColumns.FindAllStringSubmatch(demo1, -1)
	expect1 := [][]string{{"$1", "$1", "$1", "", "", "", "", "", "", ""}}
	a.Equal(expect1, res1)

	demo2 := "$2(string)"
	res2 := regexColumns.FindAllStringSubmatch(demo2, -1)
	expect2 := [][]string{{"$2(string)", "$2", "$2", "", "", "", "", "", "(string)", "string"}}
	a.Equal(expect2, res2)

	demo3 := "aggregation.key1(string)"
	res3 := regexColumns.FindAllStringSubmatch(demo3, -1)
	expect3 := [][]string{{"aggregation.key1(string)", "aggregation.key1", "", "aggregation", "aggregation", "", "", "key1", "(string)", "string"}}
	a.Equal(expect3, res3)

	preprocessor := NewPreprocessor()
	columnsCnf := map[string]string{
		"time":     "$1(time)",
		"time_utc": "$1(time_utc)",
		"name":     "aggregation.key1",
		"tags":     "static.a(int32)",
		"action":   "$2",
		"duration": "$4(int32)",
		"id":       "dynamic.id(int32)",
	}
	preprocessor.SetColumns(columnsCnf)
	preprocessingCnf := map[string]map[string]string{
		"aggregation": {
			"key1": "$2[1:3]+\"$key \"+$3[:2]",
		},
		"dynamic": {
			"id": "gen_uuid()",
		},
		"static": {
			"a": "1",
		},
	}
	preprocessor.SetPreprocessingConfig(preprocessingCnf)
	preprocessor.LoadConfig()
	idx := preprocessor.GetIndexOfFile()
	parser, _ := file_parser.DefaultParserController.GetParser("file")
	parser.SetFormat(idx)
	res4 := parser.Parse("hello world  tyltr 18")
	expect4 := map[string]string{"$1": "hello", "$2": "world", "$3": "tyltr", "$4": "18"}
	a.Equal(expect4, res4)
	expect5 := map[string]string{"key1": "or$key ty"}
	res5 := preprocessor.Do(expect4)
	a.Equal(expect5, res5)
}
