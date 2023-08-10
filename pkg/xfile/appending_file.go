package xfile

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

const AppendingFileReaderBuffSize = 1 << 4
const EventsBuffSize = 1 << 4

var (
	waitAppendingDuration = 100 * time.Millisecond
)

var (
	FatalError = errors.New("FatalError")
)

type AppendingFileAppendReader struct {
	currentCursor int64
	filename      string
	fd            *os.File
	reader        *bufio.Reader
	fileLines     chan FileLineGetter
	watcher       *fsnotify.Watcher
	fileEvents    chan EventGetter
}

func NewAppendingFileAppendReader(filename string) (XWatchReader, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	bufioReader := bufio.NewReader(fd)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	afr := &AppendingFileAppendReader{
		filename:   filename,
		fd:         fd,
		reader:     bufioReader,
		fileLines:  make(chan FileLineGetter, AppendingFileReaderBuffSize),
		watcher:    watcher,
		fileEvents: make(chan EventGetter, EventsBuffSize),
	}
	return afr, nil
}

func (afr *AppendingFileAppendReader) ReadLines(ctx context.Context) chan FileLineGetter {
	go afr.readLines(ctx)
	afr.Watch(afr.filename)
	go func() {
		for {
			select {
			case evt := <-afr.Events(ctx):
				switch evt.Operation() {
				case Chmod, Write:
					fd, err := afr.fd.Stat()
					if err != nil {
						switch {
						// file does not exist
						case os.IsNotExist(err):
							if err := afr.reWatch(); err != nil {
								line := &FileLine{
									err: FatalError,
								}
								afr.fileLines <- line
							}

						//  file does exist
						case !os.IsNotExist(err):
							line := &FileLine{
								err: FatalError,
							}
							afr.fileLines <- line
						}

					}
					if afr.currentCursor > fd.Size() {
						afr.currentCursor, err = afr.fd.Seek(0, io.SeekStart)
						if err != nil {
							line := &FileLine{
								err: FatalError,
							}
							afr.fileLines <- line
						}

						afr.reader.Reset(afr.fd)
					}
				default:
					fmt.Println("default", evt)
				}
			case <-time.After(1 * time.Second):
				fi1, err := afr.fd.Stat()
				if err != nil && !os.IsNotExist(err) {
					line := &FileLine{
						err: FatalError,
					}
					afr.fileLines <- line
				}

				fi2, err := os.Stat(afr.filename)
				if err != nil && !os.IsNotExist(err) {
					line := &FileLine{
						err: FatalError,
					}
					afr.fileLines <- line
				}

				if os.SameFile(fi1, fi2) {
					continue
				}
				if err := afr.reWatch(); err != nil {
					line := &FileLine{
						err: FatalError,
					}
					afr.fileLines <- line
				}
			}

		}
	}()
	return afr.fileLines
}

func (afr *AppendingFileAppendReader) readLines(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			close(afr.fileLines)
			afr.watcher.Close()
			afr.fd.Close()
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
			// record offset
			currentOffset, _ := afr.fd.Seek(0, io.SeekCurrent)
			afr.currentCursor = currentOffset
			line := &FileLine{
				line: b,
				err:  err,
			}
			afr.fileLines <- line
		}
	}
}

// current  offset
func (afr *AppendingFileAppendReader) CurrentCursor() int64 {
	return afr.currentCursor
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

func (afr *AppendingFileAppendReader) Watch(fileName string) error {
	return afr.watcher.Add(fileName)
}

func (afr *AppendingFileAppendReader) reWatch() error {
	afr.watcher.Remove(afr.filename)
	if err := afr.reopen(); err != nil {
		return err
	}

	return afr.watcher.Add(afr.filename)
}

func (afr *AppendingFileAppendReader) reopen() error {
	if afr.fd != nil {
		afr.fd.Close()
		afr.fd = nil
	}
	fd, err := os.Open(afr.filename)
	if err != nil {
		return err
	}

	afr.fd = fd
	afr.reader = bufio.NewReader(fd)
	return nil

}

func (afr *AppendingFileAppendReader) Events(ctx context.Context) chan EventGetter {
	go afr.watchEvents(ctx)
	return afr.fileEvents
}

func (afr *AppendingFileAppendReader) watchEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(afr.fileLines)
			afr.fd.Close()
			return
		default:
			select {
			case evt, ok := <-afr.watcher.Events:
				if !ok {
					fmt.Println("ok", ok)
					return
				}
				switch evt.Op {
				// append new data
				case fsnotify.Write:
					fileEvt := &FileEventGetter{
						evt: evt,
					}
					afr.fileEvents <- fileEvt

				default:
					fileEvt := &FileEventGetter{
						evt: evt,
					}
					fmt.Println("raw", evt)
					afr.fileEvents <- fileEvt

				}

			case err := <-afr.watcher.Errors:
				fileEvt := &FileEventGetter{err: err}
				afr.fileEvents <- fileEvt
			}
		}
	}

}
