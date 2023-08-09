package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/tylitianrui/file-clickhouse-exporter/pkg/xfile"
)

func main() {
	filename := "exmaple"
	var wg sync.WaitGroup
	var fd *os.File
	wg.Add(1)
	go func() {
		wg.Done()
		os.Remove(filename)
		fd, _ = os.Create(filename)

	}()
	wg.Wait()
	appendReader, err := xfile.NewAppendingFileAppendReader(filename)
	fmt.Println(err)
	go func() {
		for i := 0; i < 5; i++ {
			l := fmt.Sprintf("%d\n", i)
			time.Sleep(20 * time.Millisecond)
			fd.WriteString(l)
		}
	}()

	ctx, _ := context.WithCancel(context.Background())
	appendReader.Watch(filename)
	events := appendReader.Events(ctx)
	for evt := range events {
		fmt.Println("file:", evt.FileName())
		fmt.Println("operation:", evt.Operation())

	}

}
