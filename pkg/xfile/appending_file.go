package xfile

import (
	"bufio"
	"context"
	"os"
)

const AppendingFileReaderBuffSize = 1 << 4

type AppendingFileAppendReader struct {
	reader    *bufio.Reader
	fileLines chan FileLineGetter
}

func NewAppendingFileAppendReader(filename string) (XReader, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	bufioReader := bufio.NewReader(fd)

	afr := &AppendingFileAppendReader{
		reader:    bufioReader,
		fileLines: make(chan FileLineGetter, AppendingFileReaderBuffSize),
	}
	return afr, nil
}

func (afr *AppendingFileAppendReader) ReadLines(ctx context.Context) chan FileLineGetter {
	// todo
	return afr.fileLines
}
