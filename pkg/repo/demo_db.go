package repo

import (
	"context"
	"fmt"
)

type DemoRepo struct {
}

func NewDemoRepo(config ClickhouseRepoConfig) (*DemoRepo, error) {

	return &DemoRepo{}, nil

}

func (cr *DemoRepo) BatchInsert(ctx context.Context, table string, columns []string, vals [][]interface{}, async bool) error {
	fmt.Println("table:", table, " BatchInsert ok ", len(vals))
	return nil
}
