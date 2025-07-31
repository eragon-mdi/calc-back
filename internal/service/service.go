package service

import (
	"github.com/eragon-mdi/calc-back/internal/transport"
)

type service struct {
	r Repository
}

func New(r Repository) transport.Service {
	return &service{
		r: r,
	}
}
