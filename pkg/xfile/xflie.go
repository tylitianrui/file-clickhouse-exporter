package xfile

import (
	"context"
)

type XReader interface {
	ReadLines(ctx context.Context) chan FileLineGetter
}

type XWatchReader interface {
	Watch(fileName string) error
	Events() chan FileEventGetter
	XReader
}

type EventGetter interface {
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

type FileEventGetter struct{}
