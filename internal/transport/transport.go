package transport

import (
	"github.com/eragon-mdi/calc-back/internal/common/api"
	resttransport "github.com/eragon-mdi/calc-back/internal/transport/http/rest"
	"go.uber.org/zap"
)

type Service interface {
	resttransport.Service
}

func New(s Service, l *zap.SugaredLogger) api.Transport {
	return resttransport.New(s, l)
}
