package xfile

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAppendingFileAppendReader_Read(t *testing.T) {
	a := assert.New(t)
	filename := "../../test/test_appending_file"
	var wg sync.WaitGroup
	var fd *os.File
	wg.Add(1)
	go func() {
		wg.Done()
		os.Remove(filename)
		fd, _ = os.Create(filename)

	}()
	wg.Wait()
	go func() {
		for i := 0; i < 50; i++ {
			l := fmt.Sprintf("%d\n", i)
			time.Sleep(20 * time.Millisecond)
			fd.WriteString(l)
		}
	}()

	appendReader, err := NewAppendingFileAppendReader(filename)
	a.NoError(err)
	a.Implements((*XReader)(nil), appendReader)
	ctx, cancel := context.WithCancel(context.Background())
	var i int
	content := appendReader.ReadLines(ctx)
	for line := range content {
		expect := fmt.Sprintf("%d\n", i)
		a.Equal(expect, string(line.Line()))
		i++
		if i == 50 {
			cancel()
		}
	}
}

func TestAppendingFileAppendReader_Watch(t *testing.T) {
	a := assert.New(t)
	filename := "../../test/test_appending_file"
	var wg sync.WaitGroup
	var fd *os.File
	wg.Add(1)
	go func() {
		wg.Done()
		os.Remove(filename)
		fd, _ = os.Create(filename)

	}()
	wg.Wait()
	appendReader, err := NewAppendingFileAppendReader(filename)
	a.NoError(err)
	a.Implements((*XWatchReader)(nil), appendReader)
	go func() {
		for i := 0; i < 50; i++ {
			l := fmt.Sprintf("%d\n", i)
			time.Sleep(20 * time.Millisecond)
			fd.WriteString(l)
		}
	}()

	// ctx, cancel := context.WithCancel(context.Background())
	// appendReader.Watch()

}
