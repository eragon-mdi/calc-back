package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/eragon-mdi/calc-back/internal/common/configs"
)

type Storage = *sql.DB

type storage struct {
	s   Storage
	on  sync.Once
	err error
}

var store = &storage{}

func Conn(ctx context.Context, cfg *configs.Storage, driver driver, timeout time.Duration) (Storage, error) {
	store.on.Do(func() {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		store.s, store.err = conn(ctx, cfg, driver)
	})

	if store.err != nil {
		return nil, store.err
	}

	return store.s, nil
}

func conn(ctx context.Context, cfg *configs.Storage, driver driver) (Storage, error) {

	source := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLmode,
	)

	conn, err := sql.Open(driver.Name(), source)
	if err != nil {
		return nil, err
	}

	if err := conn.PingContext(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}

type driver interface {
	_mustUseWithImportedSQLDriver
	Name() string
}

type _mustUseWithImportedSQLDriver interface {
	MustUseWithImportedSQLDriver()
}
