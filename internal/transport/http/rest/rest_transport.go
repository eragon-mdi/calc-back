package resttransport

import (
	apirest "github.com/eragon-mdi/calc-back/internal/common/api/rest"
	"go.uber.org/zap"
)

type transport struct {
	s Service
	l *zap.SugaredLogger
}

func New(s Service, l *zap.SugaredLogger) apirest.Transport {
	return &transport{
		s: s,
		l: l,
	}
}
