package xfile

import (
	"bufio"
	"os"
)

type FileReader struct {
	reader *bufio.Reader
}

func NewFileReader(filename string) (*FileReader, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	bufioReader := bufio.NewReader(fd)
	fileReader := &FileReader{
		reader: bufioReader,
	}
	return fileReader, nil
}

func (fr *FileReader) ReadLine() ([]byte, error) {
	return fr.reader.ReadBytes('\n')
}
