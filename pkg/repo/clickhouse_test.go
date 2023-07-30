package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ddl = `
CREATE TABLE
    benchmark(
        id UUID NOT NULL COMMENT 'Primary KEY',
        Col1 UInt64,
        Col2 String,
        Col3 Array(UInt8),
        Col5 DateTime
    ) ENGINE = MergeTree()
ORDER BY (id) PRIMARY KEY(id) COMMENT 'test';
`

func TestNewClickhouseRepo(t *testing.T) {
	a := assert.New(t)
	config := ClickhouseRepoConfig{
		Host:     "127.0.0.1",
		Port:     9000,
		DB:       "default",
		User:     "default",
		Password: "",
	}
	conn, err := NewClickhouseRepo(config)
	a.NoError(err)
	err = conn.Ping(context.TODO())
	a.NoError(err)
	err = conn.Exec(context.TODO(), "DROP TABLE IF EXISTS benchmark")
	a.NoError(err)
	err = conn.Exec(context.TODO(), ddl)
	a.NoError(err)

	item0 := []interface{}{"hello", []uint8{1, 2, 3}}
	item1 := []interface{}{"hello2", []uint8{4, 5, 6}}
	err = conn.batchInsert(context.TODO(), "benchmark", []string{"Col2", "Col3"}, [][]interface{}{item0, item1}, false)
	a.NoError(err)
}
