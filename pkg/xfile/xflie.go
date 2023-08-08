package xfile

import (
	"bufio"
	"context"
	"os"
)

const BuffSize = 1 << 4

type XReader interface {
	ReadLines(ctx context.Context) (chan []byte, chan error)
}

type FileLineGetter interface {
	Line() []byte
	Error() error
}
type FileLine struct {
	line []byte
	err  error
}

func (fl *FileLine) Line() []byte {
	return fl.line
}

func (fl *FileLine) Error() error {
	return fl.err
}

type FileReader struct {
	reader    *bufio.Reader
	fileLines chan FileLineGetter
}

func NewFileReader(filename string) (*FileReader, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	bufioReader := bufio.NewReader(fd)
	fileReader := &FileReader{
		reader:    bufioReader,
		fileLines: make(chan FileLineGetter, BuffSize),
	}
	return fileReader, nil
}

func (fr *FileReader) readLines(ctx context.Context) {
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

func (fr *FileReader) ReadLines(ctx context.Context) chan FileLineGetter {
	go fr.readLines(ctx)
	return fr.fileLines
}
