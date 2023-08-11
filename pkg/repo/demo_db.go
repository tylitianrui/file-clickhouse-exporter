package repo

import (
	"bytes"
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

	for _, val := range vals {
		var columnVal bytes.Buffer
		for _, v := range val {
			vstr := fmt.Sprintf("%v  ", v)
			columnVal.WriteString(vstr)
		}
		fmt.Println(columnVal.String())
	}
	return nil
}
