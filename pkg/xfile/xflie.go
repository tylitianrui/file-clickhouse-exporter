package xfile

import (
	"context"

	"github.com/fsnotify/fsnotify"
)

type Operation uint32

const (
	Create Operation = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

func (op Operation) String() string {
	switch op {
	case Create:
		return "CREATE"
	case Write:
		return "WRITE"
	case Remove:
		return "REMOVE"
	case Rename:
		return "RENAME"
	case Chmod:
		return "CHMOD"
	default:
		return "UNKNOWN"
	}
}

type XReader interface {
	ReadLines(ctx context.Context) chan FileLineGetter
}

type XWatchReader interface {
	Watch(fileName string) error
	Events(ctx context.Context) chan EventGetter
	XReader
}

type EventGetter interface {
	FileName() string
	Operation() Operation
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

type FileEventGetter struct {
	evt fsnotify.Event
	err error
}

func (f *FileEventGetter) FileName() string {
	return f.evt.Name
}

func (f *FileEventGetter) Operation() Operation {
	op := Operation(f.evt.Op)
	return op

}
