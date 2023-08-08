package xfile

import (
	"bufio"
	"context"
	"os"
)

const BuffSize = 1 << 4

// StaticFileReader  reader static file.
type StaticFileReader struct {
	reader    *bufio.Reader
	fileLines chan FileLineGetter
}

func NewFileReader(filename string) (*StaticFileReader, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	bufioReader := bufio.NewReader(fd)
	fileReader := &StaticFileReader{
		reader:    bufioReader,
		fileLines: make(chan FileLineGetter, BuffSize),
	}
	return fileReader, nil
}

func (fr *StaticFileReader) readLines(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(fr.fileLines)
			return
		default:
			b, err := fr.reader.ReadBytes('\n')
			line := &FileLine{
				line: b,
				err:  err,
			}
			fr.fileLines <- line
		}
	}
}

func (fr *StaticFileReader) ReadLines(ctx context.Context) chan FileLineGetter {
	go fr.readLines(ctx)
	return fr.fileLines
}
