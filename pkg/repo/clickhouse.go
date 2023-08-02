package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
)

type ClickhouseRepoConfig struct {
	Host     string // 服务端主机
	Port     int    // 端口
	DB       string // 数据库
	User     string // 用户名
	Password string // 密码
}

type ClickhouseRepo struct {
	conn   driver.Conn
	config ClickhouseRepoConfig
}

func NewClickhouseRepo(config ClickhouseRepoConfig) (*ClickhouseRepo, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", config.Host, config.Port)},
		Auth: clickhouse.Auth{
			Database: config.DB,
			Username: config.User,
			Password: config.Password,
		},
		DialTimeout:     50 * time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		return nil, err
	}
	clickhouseRepo := &ClickhouseRepo{
		conn:   conn,
		config: config,
	}

	err = clickhouseRepo.conn.Ping(context.TODO())
	if err != nil {
		return nil, err
	}

	return clickhouseRepo, err

}

func (cr *ClickhouseRepo) Ping(ctx context.Context) error {
	return cr.conn.Ping(ctx)
}

func (cr *ClickhouseRepo) BatchInsert(ctx context.Context, table string, columns []string, vals [][]interface{}, async bool) error {
	insertSql := fmt.Sprintf("INSERT INTO %s.%s (%s)", cr.config.DB, table, strings.Join(columns, ","))
	batch, err := cr.conn.PrepareBatch(context.Background(), insertSql)
	if err != nil {
		return err
	}
	for _, val := range vals {
		err = batch.Append(val...)
		if err != nil {
			return err
		}
	}
	return batch.Send()
}

func (cr *ClickhouseRepo) Exec(ctx context.Context, sql string) error {
	return cr.conn.Exec(ctx, sql)
}
