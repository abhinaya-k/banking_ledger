package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.Logger
)

func init() {

	logConfig := zap.Config{
		OutputPaths:   []string{"stdout"},
		Encoding:      "json",
		Level:         zap.NewAtomicLevelAt(zap.DebugLevel),
		InitialFields: map[string]interface{}{"serviceName": os.Getenv("SERVICE_NAME")},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error

	Log, err = logConfig.Build()

	if err != nil {
		panic(err)
	}

}
