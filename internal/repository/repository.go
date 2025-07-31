package repository

import (
	"github.com/eragon-mdi/calc-back/internal/common/storage"
	sqlrepo "github.com/eragon-mdi/calc-back/internal/repository/sql"
	"github.com/eragon-mdi/calc-back/internal/service"
)

func New(s storage.Storage) service.Repository {
	return sqlrepo.New(s.SQL())
}
