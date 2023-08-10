package main

import (
	"context"
	"fmt"

	"github.com/tylitianrui/file-clickhouse-exporter/pkg/file_parser"
	"github.com/tylitianrui/file-clickhouse-exporter/pkg/xfile"
)

func main() {
	filename := "example"
	appendReader, err := xfile.NewAppendingFileAppendReader(filename)
	idx := []string{"$1", "$2"}
	fileParser, _ := file_parser.DefaultParserController.GetParser("file")
	fileParser.SetFormat(idx)
	fmt.Println(err)
	ctx := context.Background()
	lines := appendReader.ReadLines(ctx)
	for line := range lines {
		fmt.Println(string(line.Line()))
		res := fileParser.Parse(string(line.Line()))
		fmt.Println(res)

	}
}
