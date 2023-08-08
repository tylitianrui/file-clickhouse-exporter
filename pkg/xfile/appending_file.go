package xfile

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"time"
)

const AppendingFileReaderBuffSize = 1 << 4

var (
	waitAppendingDuration = 100 * time.Millisecond
)

type AppendingFileAppendReader struct {
	currentCursor int64
	fd            *os.File
	reader        *bufio.Reader
	fileLines     chan FileLineGetter
}

func NewAppendingFileAppendReader(filename string) (XReader, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	bufioReader := bufio.NewReader(fd)

	afr := &AppendingFileAppendReader{
		fd:        fd,
		reader:    bufioReader,
		fileLines: make(chan FileLineGetter, AppendingFileReaderBuffSize),
	}
	return afr, nil
}

func (afr *AppendingFileAppendReader) ReadLines(ctx context.Context) chan FileLineGetter {
	go afr.readLines(ctx)
	return afr.fileLines
}

func (afr *AppendingFileAppendReader) readLines(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(afr.fileLines)
			return
		default:
			b, err := afr.reader.ReadBytes('\n')
			if err != nil {
				// if we encounter EOF before a line delimiter
				// rewind cursor position, and wait for further file changes.
				if err == io.EOF {
					afr.setCursorBack(len(b))
					// todo:wait for file appending
					time.Sleep(waitAppendingDuration)
					continue
				} else {
					// other errors
					line := &FileLine{
						err: err,
					}
					afr.fileLines <- line
				}
			}
			line := &FileLine{
				line: b,
				err:  err,
			}
			afr.fileLines <- line
		}
	}
}

func (afr *AppendingFileAppendReader) setCursorBack(n int) error {
	if n < 0 {
		return errors.New("n should be positive")
	}
	if n == 0 {
		return nil
	}
	offset, err := afr.fd.Seek(-int64(n), io.SeekCurrent)
	if err != nil {
		return err
	}
	afr.currentCursor = offset
	afr.reader.Reset(afr.fd)
	return nil
}
