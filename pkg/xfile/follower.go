package xfile

import (
	"bufio"
	"io"
	"os"
)

const FILEREADERBUFSZIE = 1 << 8

type FileFollowerReader struct {
	fileReader *FileReader
	lines      chan string
	fd         *os.File
	waitOffset int64
}

func NewFileFollowerReader(filename string) (*FileFollowerReader, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	bufioReader := bufio.NewReader(fd)
	fileReader := &FileReader{
		reader: bufioReader,
	}
	fileFollowerReader := &FileFollowerReader{
		fileReader: fileReader,
	}
	return fileFollowerReader, nil
}

func (ffr *FileFollowerReader) ReadLines(line chan []byte, errch chan error) {
	for {
		for {
			b, e := ffr.fileReader.ReadLine()
			if e != nil {
				if e == io.EOF {
					l := len(b)
					ffr.waitOffset, e = ffr.fd.Seek(-int64(l), io.SeekCurrent)
					errch <- e
					reader := ffr.fileReader.Reader()

					reader.Reset(ffr.fd)

				}
			}
			line <- b
		}
	}
}
