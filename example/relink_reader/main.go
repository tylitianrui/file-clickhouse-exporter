package main

import (
	"context"
	"fmt"

	"github.com/tylitianrui/file-clickhouse-exporter/pkg/xfile"
)

func main() {
	filename := "example"
	appendReader, err := xfile.NewAppendingFileAppendReader(filename)
	fmt.Println(err)
	ctx := context.Background()
	lines := appendReader.ReadLines(ctx)
	for line := range lines {
		fmt.Println(string(line.Line()))

	}

}
