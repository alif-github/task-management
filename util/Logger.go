package util

import (
	"github.com/alif-github/task-management/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger

type LoggerModel struct {
	RequestID   string `json:"request_id"`
	Class       string `json:"class"`
	Application string `json:"application"`
	Version     string `json:"version"`
	Code        int    `json:"code"`
	Message     string `json:"message"`
}

func GenerateLogModel(version string, application string) (output LoggerModel) {
	output.RequestID = "-"
	output.Class = "-"
	output.Application = application
	output.Version = version
	output.Code = 0
	output.Message = "-"
	return
}

func DefaultGenerateLogModel(code int, msg string) (output LoggerModel) {
	output.RequestID = "-"
	output.Class = "-"
	output.Application = config.ApplicationConfiguration.GetServer().Application
	output.Version = config.ApplicationConfiguration.GetServer().Version
	output.Code = code
	output.Message = msg
	return
}

func (object LoggerModel) LoggerZapFieldObject() (output []zap.Field) {
	output = append(output, zap.String("requestID", object.RequestID))
	output = append(output, zap.String("class", object.Class))
	output = append(output, zap.String("application", object.Application))
	output = append(output, zap.String("version", object.Version))
	output = append(output, zap.Int("code", object.Code))
	output = append(output, zap.String("message", object.Message))
	return
}

func ConfigZap(logFile []string) {
	var (
		cfg zap.Config
		err error
	)

	cfg = zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:    "json",
		OutputPaths: logFile,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "timestamp",
			EncodeTime:  zapcore.RFC3339NanoTimeEncoder,
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
	}

	Logger, err = cfg.Build()
	if err != nil {
		os.Exit(3)
	}

	return
}

func LogInfo(data []zap.Field) {
	Logger.Info("", data...)
}

func LogError(data []zap.Field) {
	Logger.Error("", data...)
}
