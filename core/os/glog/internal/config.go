package internal

import (
	"errors"
	"github.com/Ravior/gserver/core/util/gconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

func NewSugaredLogger(config zapcore.EncoderConfig, defaultWriter zapcore.WriteSyncer, errorWriter zapcore.WriteSyncer, level zapcore.Level, opts ...zap.Option) *zap.SugaredLogger {
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			zapcore.AddSync(defaultWriter),
			zap.NewAtomicLevelAt(level),
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			zapcore.AddSync(errorWriter),
			zap.NewAtomicLevelAt(zapcore.ErrorLevel),
		),
	)

	zapLogger := zap.New(core, opts...)

	return zapLogger.Sugar()
}

func NewConsoleLogConfig() *gconfig.LogConfig {
	config := &gconfig.LogConfig{
		Level:           "debug",
		StackLevel:      "errors",
		EnableWriteFile: false,
		EnableConsole:   true,
		FilePath:        "",
		MaxAge:          0,
		TimeFormat:      "2006-01-02 15:04:05.000", //2006-01-02 15:04:05.000
		PrintCaller:     true,
	}
	return config
}

func GetLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "errors":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.FatalLevel
	}
}

// GetLogConfig 读取日志配置
func GetLogConfig(name string) (*gconfig.LogConfig, error) {
	if logConfig, ok := gconfig.Global.Log[name]; ok {
		logConfig.FilePath = gconfig.Global.BasePath + logConfig.FilePath
		return logConfig, nil
	}
	return nil, errors.New("日志配置不存在")
}
