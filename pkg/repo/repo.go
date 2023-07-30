package repo

import (
	"context"
)

type Repo interface {
	InsertM(ctx context.Context, columns []string, val [][]interface{}) error
}
