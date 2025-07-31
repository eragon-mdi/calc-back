package logging

import (
	"github.com/go-faster/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	tsKey         = "timestamp"
	ErrParseLevel = "error ParseAtomicLevel"
	ErrBuild      = "error logConfig.Build"
)

func NewLogger(level, encoding, output, messageKey string) (*zap.SugaredLogger, error) {
	logLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, errors.Wrapf(err, "%s %s", ErrParseLevel, level)
	}

	logger, err := zap.Config{
		Level:       logLevel,
		Encoding:    encoding,
		OutputPaths: []string{output},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: messageKey,
			TimeKey:    tsKey,
			EncodeTime: zapcore.RFC3339NanoTimeEncoder,
		},
		DisableStacktrace: true,
	}.Build()

	if err != nil {
		return nil, errors.Wrap(err, ErrBuild)
	}

	return logger.Sugar(), nil
}
