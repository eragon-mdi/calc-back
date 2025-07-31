package sqlrepo

import (
	"github.com/eragon-mdi/calc-back/internal/service"
	sqlstore "github.com/eragon-mdi/calc-back/pkg/storage/sql"
)

type sqlRepo struct {
	s sqlstore.Storage
}

func New(s sqlstore.Storage) service.Repository {
	return &sqlRepo{
		s: s,
	}
}
