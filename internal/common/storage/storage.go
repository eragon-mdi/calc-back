package storage

import (
	"context"
	"time"

	"github.com/eragon-mdi/calc-back/internal/common/configs"
	pgdriver "github.com/eragon-mdi/calc-back/pkg/storage/drivers/postgres"
	sqlstore "github.com/eragon-mdi/calc-back/pkg/storage/sql"
	"github.com/go-faster/errors"
)

const (
	ErrConnectSqlDB    = "Failed to connect to db"
	ErrDisconnectSqlDB = "Failed to disconnect sql-db"
	ConnTimeoutDefault = time.Minute
)

type Storage interface {
	SQL() sqlstore.Storage
	GracefulShutdown() error
}

type storage struct {
	sqlStore sqlstore.Storage
}

func Conn(cfg *configs.Storage, timeout time.Duration) (Storage, error) {
	sql, err := sqlstore.Conn(context.Background(), cfg, pgdriver.Postgres{}, timeout)
	if err != nil {
		return nil, errors.Wrap(err, ErrConnectSqlDB)
	}

	return &storage{
		sqlStore: sql,
	}, nil
}

func (s storage) SQL() sqlstore.Storage {
	return s.sqlStore
}

func (s storage) GracefulShutdown() error {
	if err := s.sqlStore.Close(); err != nil {
		return errors.Wrap(err, ErrDisconnectSqlDB)
	}

	return nil
}
